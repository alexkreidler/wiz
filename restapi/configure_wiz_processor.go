// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/alexkreidler/wiz/restapi/operations"
	"github.com/alexkreidler/wiz/restapi/operations/processors"
)

//go:generate swagger generate server --target ../../wiz --name WizProcessor --spec ../api/swagger.json

func configureFlags(api *operations.WizProcessorAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.WizProcessorAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	if api.ProcessorsAddDataHandler == nil {
		api.ProcessorsAddDataHandler = processors.AddDataHandlerFunc(func(params processors.AddDataParams) middleware.Responder {
			return middleware.NotImplemented("operation processors.AddData has not yet been implemented")
		})
	}
	if api.ProcessorsGetAllProcessorsHandler == nil {
		api.ProcessorsGetAllProcessorsHandler = processors.GetAllProcessorsHandlerFunc(func(params processors.GetAllProcessorsParams) middleware.Responder {
			return middleware.NotImplemented("operation processors.GetAllProcessors has not yet been implemented")
		})
	}
	if api.ProcessorsGetAllRunsHandler == nil {
		api.ProcessorsGetAllRunsHandler = processors.GetAllRunsHandlerFunc(func(params processors.GetAllRunsParams) middleware.Responder {
			return middleware.NotImplemented("operation processors.GetAllRuns has not yet been implemented")
		})
	}
	if api.ProcessorsGetConfigHandler == nil {
		api.ProcessorsGetConfigHandler = processors.GetConfigHandlerFunc(func(params processors.GetConfigParams) middleware.Responder {
			return middleware.NotImplemented("operation processors.GetConfig has not yet been implemented")
		})
	}
	if api.ProcessorsGetDataHandler == nil {
		api.ProcessorsGetDataHandler = processors.GetDataHandlerFunc(func(params processors.GetDataParams) middleware.Responder {
			return middleware.NotImplemented("operation processors.GetData has not yet been implemented")
		})
	}
	if api.ProcessorsGetInputChunkHandler == nil {
		api.ProcessorsGetInputChunkHandler = processors.GetInputChunkHandlerFunc(func(params processors.GetInputChunkParams) middleware.Responder {
			return middleware.NotImplemented("operation processors.GetInputChunk has not yet been implemented")
		})
	}
	if api.ProcessorsGetOutputChunkHandler == nil {
		api.ProcessorsGetOutputChunkHandler = processors.GetOutputChunkHandlerFunc(func(params processors.GetOutputChunkParams) middleware.Responder {
			return middleware.NotImplemented("operation processors.GetOutputChunk has not yet been implemented")
		})
	}
	if api.ProcessorsGetProcessorHandler == nil {
		api.ProcessorsGetProcessorHandler = processors.GetProcessorHandlerFunc(func(params processors.GetProcessorParams) middleware.Responder {
			return middleware.NotImplemented("operation processors.GetProcessor has not yet been implemented")
		})
	}
	if api.ProcessorsGetRunHandler == nil {
		api.ProcessorsGetRunHandler = processors.GetRunHandlerFunc(func(params processors.GetRunParams) middleware.Responder {
			return middleware.NotImplemented("operation processors.GetRun has not yet been implemented")
		})
	}
	if api.ProcessorsUpdateConfigHandler == nil {
		api.ProcessorsUpdateConfigHandler = processors.UpdateConfigHandlerFunc(func(params processors.UpdateConfigParams) middleware.Responder {
			return middleware.NotImplemented("operation processors.UpdateConfig has not yet been implemented")
		})
	}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
