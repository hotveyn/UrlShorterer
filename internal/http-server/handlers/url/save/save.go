package save

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
	"net/http"
	"server/internal/lib/api/response"
	"server/internal/lib/logger/sl"
	"server/internal/lib/random"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			fmt.Println(err)
			//log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)
			log.Error("Validation failed", sl.Err(validateErr))

			render.JSON(w, r, response.ValidationError(validateErr))

			return
		}

		const aliasLength = 6
		alias := req.Alias
		if alias == "" {
			random.String(&alias, aliasLength)
		}

		//result, err := urlSaver.SaveURL(req.URL, alias)
		//if errors.Is(err, storage.ErrURLExists) {
		//	log.Info("url already exists", slog.String("url", req.URL))
		//
		//	render.JSON(w, r, response.Error("URL already exists"))
		//
		//	return
		//}
		//
		//log.Info("url added", slog.Int64("id", result))

		render.JSON(w, r, Response{
			Response: response.OK(),
			Alias:    alias,
		})
	}
}
