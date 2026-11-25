// package main indicates that this is an executable Go program.
package main

import (
	"context"       // Manages cancellation, deadlines, and timeouts.
	"encoding/json" // Encodes and decodes JSON data.
	"errors"        // Provides functions for handling and comparing errors.
	"log/slog"      // Provides structured logging.
	"net/http"      // Provides HTTP client and server functionality.
	"os"            // Provides operating system functionality.
	"os/signal"     // Handles operating system signals.
	"syscall"       // Provides low-level system constants and operations.
	"time"          // Handles time and duration.
)

// main() is the entry point of a Go executable program.
func main() {

	// Create a structured logger that writes logs to the terminal.
	// os.Stdout represents the standard output.
	// nil means we use the default handler configuration.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Create an HTTP request multiplexer (router).
	// It routes incoming HTTP requests to the appropriate handlers.
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {

		logger.Info("health check", "remote", r.RemoteAddr)

		// Set the response header to indicate that the response is JSON.
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(
			map[string]string{"status": "ok"},
		); err != nil {

			// Log an error if encoding or writing the response fails.
			logger.Error("write health response", "error", err)
		}
	})

	// Create an HTTP server using a struct literal.
	// The & operator returns a pointer to the server struct.
	srv := &http.Server{
		Addr: ":8082",

		// Use the router created above to handle incoming requests.
		Handler: mux,

		// Allow a maximum of 5 seconds to read an entire request, including its body.
		ReadTimeout: 5 * time.Second,
	}

	// Create a context that is automatically cancelled when the
	// application receives SIGINT or SIGTERM.
	//
	// context.Background(): Creates a root context.
	// SIGINT: Interrupt signal, usually triggered by Ctrl+C.
	// SIGTERM: Signal requesting the application to terminate.
	//
	// ctx: The context used to monitor cancellation.
	// stop: A function that releases signal-handling resources.
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	// defer schedules a function to execute when main() returns.
	// It ensures that signal-handling resources are released.
	defer stop()

	// Start the HTTP server in a goroutine.
	//
	// A goroutine is a lightweight concurrent execution unit in Go.
	// The "go" keyword allows this function to run concurrently
	// without blocking the main function.
	go func() {

		// Log the address on which the server is starting.
		logger.Info("starting server", "addr", srv.Addr)

		if err := srv.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {

			// Log unexpected server errors.
			logger.Error("server shutdown", "error", err)

			// Exit the program with status code 1, indicating failure.
			os.Exit(1)
		}

	}() // Immediately execute the anonymous function as a goroutine.

	// Block the main goroutine until the context is cancelled.
	//
	// ctx.Done() returns a channel that closes when cancellation occurs.
	// The <- operator receives from a channel.
	//
	// This prevents main() from exiting while the server is running.
	<-ctx.Done()

	logger.Info("shutdown signal received")

	// Create a separate context with a 5-second timeout.
	// This gives the server time to finish active HTTP requests.
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	// Release the timeout context's resources when main() returns.
	defer cancel()

	// Gracefully shut down the HTTP server.
	// Shutdown stops accepting new connections and waits for
	// active requests to finish, subject to the context timeout.
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
		return
	}

	logger.Info("server stopped successfully")
}