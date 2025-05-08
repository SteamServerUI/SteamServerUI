package web

import (
	"net/http"
	"net/http/pprof"
	"os"
	"sync"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v6/src/config"
	"github.com/SteamServerUI/SteamServerUI/v6/src/logger"
	"github.com/SteamServerUI/SteamServerUI/v6/src/misc"
	"github.com/SteamServerUI/SteamServerUI/v6/src/security"
)

func StartWebServer(wg *sync.WaitGroup) {
	logger.Web.Info("Starting API routes...")

	// Set up routes
	mux, protectedMux := SetupRoutes()

	// Apply middleware to protected routes
	mux.Handle("/", AuthMiddleware(protectedMux)) // Wrap protected routes under root

	backendEndpoint := config.GetBackendEndpointIP() + ":" + config.GetBackendEndpointPort()
	pprofEndpoint := config.GetBackendEndpointIP() + ":6060"
	backendEndpointUrl := "https://" + backendEndpoint

	// Start HTTP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Web.Info("Starting the HTTP server...")
		// Ensure TLS certs are ready
		if err := security.EnsureTLSCerts(); err != nil {
			logger.Web.Error("Error setting up TLS certificates: " + err.Error())
			exitServerWithDelay("TLS Certificate Error.", 20)
		}
		misc.PrintStartupMessage(backendEndpointUrl)
		if config.GetIsFirstTimeSetup() {
			misc.PrintFirstTimeSetupMessage()
		}
		logger.Core.Info("Ready to run your server!")
		logger.Core.Info("🙏Thank you for using SSUI!")
		err := http.ListenAndServeTLS(backendEndpoint, config.GetTLSCertPath(), config.GetTLSKeyPath(), mux)
		if err != nil {
			logger.Web.Error("Error starting HTTPS server: " + err.Error())
			exitServerWithDelay("Error starting HTTPS server.", 20)
		}
	}()

	// Start the pprof server if debug mode is enabled (HTTP/1.1)
	if config.GetIsDebugMode() && config.GetLogLevel() < 20 { // if debug mode is enabled and log level is lower than 20 (if this triggers LogLevel is probably 10 and probably debug, but who knows), start pprof server
		wg.Add(1)
		go func() {
			defer wg.Done()
			pprofMux := http.NewServeMux()
			// Register pprof handler
			pprofMux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
			logger.Web.Warn("⚠️Starting pprof server on" + pprofEndpoint + "/debug/pprof")
			pprofIPandPort := pprofEndpoint
			err := http.ListenAndServe(pprofIPandPort, pprofMux)
			if err != nil {
				logger.Web.Error("Error starting pprof server: " + err.Error())
			}
		}()
	}

	// Wait for both servers to be running
	wg.Wait()

}

func exitServerWithDelay(message string, timeToExit int) {
	logger.Web.Error("Not healthy, I have to go...: " + message)
	ttE := time.Duration(timeToExit) * time.Second
	time.Sleep(ttE)
	os.Exit(1)
}
