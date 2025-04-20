package web

import (
	"net/http"
	"net/http/pprof"
	"sync"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/security"
)

func StartWebServer(wg *sync.WaitGroup) {
	logger.Web.Info("Starting API routes...")

	// Set up routes
	mux, protectedMux := SetupRoutes()

	// Apply middleware to protected routes
	mux.Handle("/", AuthMiddleware(protectedMux)) // Wrap protected routes under root

	// Start HTTP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Web.Info("Starting the HTTP server on port 8443...")
		logger.Web.Info("UI available at: https://0.0.0.0:8443 or https://localhost:8443")
		if config.IsFirstTimeSetup {
			logger.Web.Error("For first-time setup, visit the UI to configure a user or skip authentication.")
			logger.Web.Warn("Fill the Username and Password fields, then click Register User and when done Finalize Setup.")
			logger.Web.Warn("For more details, check the GitHub Wiki: https://github.com/JacksonTheMaster/StationeersServerUI/v5/wiki")
		}
		// Ensure TLS certs are ready
		if err := security.EnsureTLSCerts(); err != nil {
			logger.Web.Error("Error setting up TLS certificates: " + err.Error())
			//os.Exit(1)
		}
		err := http.ListenAndServeTLS("0.0.0.0:8443", config.TLSCertPath, config.TLSKeyPath, mux)
		if err != nil {
			logger.Web.Error("Error starting HTTPS server: " + err.Error())
		}
	}()

	// Start the pprof server if debug mode is enabled (HTTP/1.1)
	if config.IsDebugMode && config.LogLevel < 20 { // if debug mode is enabled and log level is lower than 20 (if this triggers LogLevel is probably 10 and probably debug, but who knows), start pprof server
		wg.Add(1)
		go func() {
			defer wg.Done()
			pprofMux := http.NewServeMux()
			// Register pprof handler
			pprofMux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
			logger.Web.Warn("⚠️Starting pprof server on :6060/debug/pprof")
			err := http.ListenAndServe("0.0.0.0:6060", pprofMux)
			if err != nil {
				logger.Web.Error("Error starting pprof server: " + err.Error())
			}
		}()
	}

	// Wait for both servers to be running
	wg.Wait()

}
