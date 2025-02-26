package v1

import (
	"net/http"

	"github.com/DillonEnge/thunk/internal/api"
	"github.com/DillonEnge/thunk/templates"
)

func HandleChat() api.HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) *api.ApiError {
		templates.Chat(nil).Render(r.Context(), w)

		return nil
	}
}
