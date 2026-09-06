package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	MsgInvalidReqPayload     = "invalid request payload"
	MsgEmptyBody             = "request body must not be empty"
	MsgSyntaxErr             = "request body contains badly-formed JSON (at position %d)"
	MsgUnexpectedEOF         = "request body contains badly-formed JSON"
	MsgWrongFieldType        = "request body contains an invalid value for the %q field (at position %d)"
	MsgUnknownField          = "request body contains unknown field %s"
	MsgMultipleReqBody       = "request body must only contain a single JSON object"
	MsgMaxBodyBytes          = "request body must not be larger than %d bytes"
	MsgUnexpectedContentType = "content-Type header is not application/json"
	MsgJSONError             = "malformed JSON payload"
	PrefixUnknownFieldErr    = "json: unknown field "
	MaxBodyBytes             = 1048576
)

type ErrorResponse struct {
	Message string            `json:"message"`
	Status  int               `json:"status"`
	Errors  map[string]string `json:"errors,omitempty"`
}

func NewErrorResponse(message string, status int, errors map[string]string) *ErrorResponse {
	return &ErrorResponse{
		Message: message,
		Status:  status,
		Errors:  errors,
	}
}

type Validator interface {
	Valid() (problems map[string]string)
}

func DecodeValid[T Validator](w http.ResponseWriter, r *http.Request) (T, *ErrorResponse) {
	var v T
	if errResp := contentTypeValid(r); errResp != nil {
		return v, errResp
	}

	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&v); err != nil {
		errResp := mapJSONError(err)
		return v, &errResp
	}

	err := dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		errResp := ErrorResponse{Message: MsgMultipleReqBody, Status: http.StatusBadRequest}
		return v, &errResp
	}

	if problems := v.Valid(); len(problems) > 0 {
		errResp := ErrorResponse{
			Message: MsgInvalidReqPayload,
			Status:  http.StatusUnprocessableEntity,
			Errors:  problems,
		}
		return v, &errResp
	}
	return v, nil
}

func contentTypeValid(r *http.Request) *ErrorResponse {
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != "application/json" {
			return &ErrorResponse{
				Message: MsgUnexpectedContentType,
				Status:  http.StatusUnsupportedMediaType,
			}
		}
	}
	return nil
}

func mapJSONError(err error) ErrorResponse {
	var syntaxErr *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError
	var maxBytesError *http.MaxBytesError
	switch {
	case errors.Is(err, io.EOF):
		return ErrorResponse{
			Message: MsgEmptyBody,
			Status:  http.StatusBadRequest,
		}

	case errors.Is(err, io.ErrUnexpectedEOF):
		return ErrorResponse{
			Message: MsgUnexpectedEOF,
			Status:  http.StatusBadRequest,
		}

	case errors.As(err, &syntaxErr):
		msg := fmt.Sprintf(MsgSyntaxErr, syntaxErr.Offset)
		return ErrorResponse{
			Message: msg,
			Status:  http.StatusBadRequest,
		}

	case errors.As(err, &unmarshalTypeError):
		msg := fmt.Sprintf(MsgWrongFieldType, unmarshalTypeError.Field, unmarshalTypeError.Offset)
		return ErrorResponse{
			Message: msg,
			Status:  http.StatusBadRequest,
		}

	case strings.HasPrefix(err.Error(), PrefixUnknownFieldErr):
		fieldName := strings.TrimPrefix(err.Error(), PrefixUnknownFieldErr)
		msg := fmt.Sprintf(MsgUnknownField, fieldName)
		return ErrorResponse{
			Message: msg,
			Status:  http.StatusBadRequest,
		}

	case errors.As(err, &maxBytesError):
		msg := fmt.Sprintf(MsgMaxBodyBytes, maxBytesError.Limit)
		return ErrorResponse{
			Message: msg,
			Status:  http.StatusRequestEntityTooLarge,
		}
	default:
		return ErrorResponse{
			Message: MsgJSONError,
			Status:  http.StatusInternalServerError,
		}
	}
}
