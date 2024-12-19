package v1

import (
	"context"
	"net/http"

	"github.com/DillonEnge/thunk/internal/api"
	"github.com/DillonEnge/thunk/templates"
)

func HandleBase() api.HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) *api.ApiError {

		templates.Base().Render(context.Background(), w)

		return nil
	}
}
