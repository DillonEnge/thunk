package server

import (
	"context"
	"errors"
	"fmt"

	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/DillonEnge/thunk/internal/api"
	"github.com/DillonEnge/thunk/internal/api/middleware"
	v1 "github.com/DillonEnge/thunk/internal/api/v1"
	"github.com/DillonEnge/thunk/internal/ollama"
)

func Start(address string, config *api.Config) func(context.Context) error {
	ollamaClient := ollama.NewClient(config.Ollama.Domain, nil)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.Handle(
		"GET /static/",
		middleware.NoCache(
			http.StripPrefix("/static/",
				http.FileServer(http.Dir("./templates/static")),
			),
		),
	)

	mux.HandleFunc("GET /loader", makeH(v1.HandleLoader()))

	mux.HandleFunc("GET /", makeH(v1.HandleBase()))

	mux.HandleFunc("GET /chat", makeH(v1.HandleChat()))

	// mux.HandleFunc("GET /completion", makeH(v1.HandleCompletion(ollamaClient)))

	mux.HandleFunc("GET /ws/messages", makeH(v1.HandleMessagesWS(ollamaClient)))

	h := middleware.NewHandlerWithMiddleware(
		mux,
		middleware.Logger,
	)

	s := &http.Server{
		Addr:              address,
		Handler:           h,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		slog.Info("Listening...", "address", address)
		err := s.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	return s.Shutdown
}

func Service(ctx context.Context, config *api.Config) (func(), error) {
	shutdown := Start(fmt.Sprintf(":%d", config.Port), config)

	stopService := func() {
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		if err := shutdown(ctx); err != nil {
			panic(err)
		}
	}

	return stopService, nil
}

func makeH(h api.HandlerFuncWithError) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			w.WriteHeader(err.Status)
			errJSON := fmt.Sprintf(
				`{"error": "%s"}`,
				strings.ReplaceAll(err.Error(), `"`, `'`),
			)
			w.Write([]byte(errJSON))
		}
	}
}
