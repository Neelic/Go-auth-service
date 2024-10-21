package save

import (
	"errors"
	resp "go-url-service/internal/lib/api/response"
	"go-url-service/internal/lib/logger/sl"
	"go-url-service/internal/lib/random"
	"go-url-service/internal/storage"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URLToSave string `json:"url" validate:"required,url"`
	Alias     string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"error,omitempty"`
}

const aliasLength = 6

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.Decode(r, &req)
		if err != nil {
			log.Error("Failed to decode request", sl.Err(err))
			render.JSON(w, r, resp.Error("Failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))
		if err := validator.New().Struct(req); err != nil {
			log.Error("Failed to validate request", sl.Err(err))
			render.JSON(w, r, resp.Error("Failed to validate request"))
			render.JSON(w, r, resp.ValidationError(err.(validator.ValidationErrors)))
			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		id, err := urlSaver.SaveURL(req.URLToSave, alias)
		if errors.Is(err, storage.ErrUrlExist) {
			log.Error("url already exist", sl.Err(err))
			render.JSON(w, r, resp.Error("url already exist"))
			return
		}

		if err != nil {
			log.Error("Failed to save url", sl.Err(err))
			render.JSON(w, r, resp.Error("Failed to save url"))
			return
		}

		log.Info("url added", slog.Int64("id", id))
		render.JSON(w, r, Response{
			Response: resp.Ok(),
			Alias:    alias,
		})
	}
}
