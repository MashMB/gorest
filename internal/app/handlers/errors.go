package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/mashmb/gorest/internal/app/types"
)

func dtoToJson(dto types.ApiErrorDto) []byte {
	jsonBytes, _ := json.Marshal(dto)

	return jsonBytes
}

func HandleError(err error, res http.ResponseWriter) {
	res.Header().Set("Content-Type", "application/json")
	var apiErr *types.ApiError

	if errors.As(err, &apiErr) {
		slog.Error("Handling expected API error", "err", err.Error())

		switch {
		case apiErr.Status == 401:
			dto := types.NewApiErrorDto("auth.unauthorized", "Unauthorized")
			res.WriteHeader(apiErr.Status)
			res.Write(dtoToJson(dto))
			return
		case apiErr.Status == 400:
			dto := types.NewApiErrorDto("err.validation", "Validation Failed", apiErr.Details...)
			res.WriteHeader(apiErr.Status)
			res.Write(dtoToJson(dto))
			return
		case apiErr.Status == 406:
			dto := types.NewApiErrorDto("err.not-acceptable", "Not Acceptable")
			res.WriteHeader(apiErr.Status)
			res.Write(dtoToJson(dto))
			return
		}
	}

	slog.Error("Handling unexpected API error", "err", err.Error())
	errDto := types.NewApiErrorDto("err.internal", "Internal Server Error")
	res.WriteHeader(http.StatusInternalServerError)
	res.Write(dtoToJson(errDto))
}
