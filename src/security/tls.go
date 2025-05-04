package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v6/src/config"
	"github.com/SteamServerUI/SteamServerUI/v6/src/logger"
)

// EnsureTLSCerts ensures TLS certificates exist and are valid at config.GetTLSCertPath() and config.GetTLSKeyPath(), generating self-signed ones if needed.
func EnsureTLSCerts() error {
	// Check if cert and key files exist
	certExists := fileExists(config.GetTLSCertPath())
	keyExists := fileExists(config.GetTLSKeyPath())
	tlsDir := config.GetUIModFolder() + "tls/"

	if _, err := os.Stat(tlsDir); os.IsNotExist(err) {
		if err := os.MkdirAll(tlsDir, os.ModePerm); err != nil {
			logger.Security.Error("Failed to create directory " + tlsDir + ": " + err.Error())
			return nil
		}
	}

	if certExists && keyExists {
		// Load and check if the cert is still valid
		certData, err := os.ReadFile(config.GetTLSCertPath())
		if err != nil {
			return fmt.Errorf("failed to read cert file: %v", err)
		}
		certBlock, _ := pem.Decode(certData)
		if certBlock == nil {
			return fmt.Errorf("failed to decode PEM certificate")
		}
		cert, err := x509.ParseCertificate(certBlock.Bytes)
		if err != nil {
			return fmt.Errorf("failed to parse certificate: %v", err)
		}

		// Check if expired or near expiry (within 10 days of 90-day validity)
		if time.Now().After(cert.NotAfter) || time.Now().Add(10*24*time.Hour).After(cert.NotAfter) {
			logger.Security.Warn("Certificate expired or near expiry, regenerating...")
		} else {
			// Cert is valid, we’re done
			return nil
		}
	}

	// Generate a new self-signed cert if files are missing or expired
	if err := generateSelfSignedCert(); err != nil {
		return fmt.Errorf("failed to generate self-signed cert: %v", err)
	}

	logger.Security.Info("Generated new self-signed TLS certificates at " + config.GetTLSCertPath() + " and " + config.GetTLSKeyPath())
	return nil
}

// generateSelfSignedCert creates a self-signed certificate and key pair at config.GetTLSCertPath() and config.GetTLSKeyPath().
func generateSelfSignedCert() error {
	dir := filepath.Dir(config.GetTLSCertPath())
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", dir, err)
	}
	// Generate a private key
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// Create a certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"SteamServerUI"},
			CommonName:   "localhost",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(90 * 24 * time.Hour), // 90 days
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost", "0.0.0.0"},
	}

	// Create the certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return err
	}

	// Write the certificate to file
	certFile, err := os.Create(config.GetTLSCertPath())
	if err != nil {
		return err
	}
	defer certFile.Close()
	pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	// Write the private key to file
	keyFile, err := os.Create(config.GetTLSKeyPath())
	if err != nil {
		return err
	}
	defer keyFile.Close()
	pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	return nil
}

// fileExists checks if a file exists at the given path.
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
