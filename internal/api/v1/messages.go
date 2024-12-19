package v1

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/DillonEnge/thunk/internal/api"
	"github.com/DillonEnge/thunk/internal/ollama"
	"github.com/DillonEnge/thunk/templates"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/gofrs/uuid/v5"
)

type PostMessageParams struct {
	Message string `json:"message"`
}

func HandleMessagesWS(ollamaClient *ollama.Client) api.HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) *api.ApiError {
		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			return &api.ApiError{
				Status: http.StatusInternalServerError,
				Err:    err,
			}
		}
		defer c.CloseNow()

		// Set the context as needed. Use of r.Context() is not recommended
		// to avoid surprising behavior (see http.Hijacker).
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute*30)
		defer cancel()

		messages := []ollama.Message{}

		for {
			var params PostMessageParams
			err = wsjson.Read(ctx, c, &params)
			if err != nil {
				slog.Error("error reading json from ws", "err", err)
				break
			}

			if params.Message == "" {
				slog.Warn("encountered blank message, skipping...", "params", params)
				continue
			}

			slog.Info("recieved message from ws", "msg", params)

			messages = append(messages, ollama.Message{
				Role:    "user",
				Content: params.Message,
			})

			d, _, err := buildMessageOOB(ctx, params.Message, "end", false)
			if err != nil {
				slog.Warn("failed to build message OOB")
				continue
			}

			c.Write(ctx, websocket.MessageText, d)

			d, id, err := buildMessageOOB(ctx, "", "start", true)
			if err != nil {
				slog.Warn("failed to build message OOB")
				continue
			}

			c.Write(ctx, websocket.MessageText, d)

			resp, err := ollamaClient.ChatCompletion("qwen2.5:7b", messages)
			if err != nil {
				slog.Error("error getting completion", "err", err)
				break
			}

			var m string
			s := bufio.NewScanner(resp.Body)

			for s.Scan() {
				var r ollama.ChatCompletionResponse
				json.Unmarshal(s.Bytes(), &r)

				m += r.Message.Content

				slog.Info("Received message from ollama", "message", r.Message.Content)

				var buf bytes.Buffer
				templates.MessageContentOOB(m, id).Render(ctx, &buf)

				d, msgErr := io.ReadAll(&buf)
				if msgErr != nil {
					slog.Error("failed to read from message content oob buf", "err", err)
					continue
				}

				c.Write(ctx, websocket.MessageText, d)
			}

			messages = append(messages, ollama.Message{
				Role:    "assistant",
				Content: m,
			})
		}

		c.Close(websocket.StatusNormalClosure, "")

		slog.Error("encountered err", "error", err)

		return nil
	}
}

func buildMessageOOB(ctx context.Context, message string, position string, withID bool) ([]byte, string, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return nil, "", err
	}

	id := ""

	if withID {
		id = uid.String()
	}

	var buf bytes.Buffer
	templates.MessageOOB(message, position, id).Render(ctx, &buf)

	d, err := io.ReadAll(&buf)
	if err != nil {
		slog.Error("failed to read from message oob buf", "err", err)
		return nil, "", err
	}

	return d, uid.String(), nil
}
