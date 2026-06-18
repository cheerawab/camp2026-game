package httpx

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	chimw "github.com/go-chi/chi/v5/middleware"
)

type loggerContextKey struct{}

type Error struct {
	Status  int
	Detail  string
	Code    string
	Cause   error
	Details []ErrorDetail
}

func (e *Error) Error() string {
	return e.Detail
}

func (e *Error) Unwrap() error {
	return e.Cause
}

type ProblemDetails struct {
	Type      string        `json:"type,omitempty"`
	Title     string        `json:"title,omitempty"`
	Status    int           `json:"status,omitempty"`
	Detail    string        `json:"detail,omitempty"`
	Code      string        `json:"code,omitempty"`
	RequestID string        `json:"requestId,omitempty"`
	Instance  string        `json:"instance,omitempty"`
	Errors    []ErrorDetail `json:"errors,omitempty"`
}

type ErrorDetail struct {
	Location string `json:"location,omitempty"`
	Message  string `json:"message,omitempty"`
}

func NewError(status int, detail string) error {
	return &Error{Status: status, Detail: detail}
}

func InternalServerError(detail string, code string, cause error) error {
	return &Error{Status: http.StatusInternalServerError, Detail: detail, Code: code, Cause: cause}
}

func BadRequest(detail string) error {
	return NewError(http.StatusBadRequest, detail)
}

func NotFound(detail string) error {
	return NewError(http.StatusNotFound, detail)
}

func UnprocessableEntity(detail string, details ...ErrorDetail) error {
	return &Error{Status: http.StatusUnprocessableEntity, Detail: detail, Details: details}
}

func ServiceUnavailable(detail string) error {
	return NewError(http.StatusServiceUnavailable, detail)
}

func WithLogger(ctx context.Context, log *slog.Logger) context.Context {
	if log == nil {
		return ctx
	}
	return context.WithValue(ctx, loggerContextKey{}, log)
}

func WriteJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if body == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(body)
}

func WriteProblem(w http.ResponseWriter, r *http.Request, err error) {
	status := http.StatusInternalServerError
	detail := http.StatusText(status)
	code := ""

	var httpErr *Error
	if errors.As(err, &httpErr) {
		status = httpErr.Status
		detail = httpErr.Detail
		code = httpErr.Code
	}

	requestID := chimw.GetReqID(r.Context())
	if requestID != "" {
		w.Header().Set("X-Request-ID", requestID)
	}
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)

	body := ProblemDetails{
		Status: status,
		Title:  http.StatusText(status),
		Detail: detail,
	}
	if status >= http.StatusInternalServerError {
		if code == "" {
			code = defaultServerErrorCode(status)
		}
		body.Code = code
		body.RequestID = requestID
	}
	if httpErr != nil && len(httpErr.Details) > 0 {
		body.Errors = httpErr.Details
	}

	logProblem(r, status, detail, code, requestID, err, httpErr)

	_ = json.NewEncoder(w).Encode(body)
}

func logProblem(
	r *http.Request,
	status int,
	detail string,
	code string,
	requestID string,
	err error,
	httpErr *Error,
) {
	if status < http.StatusInternalServerError {
		return
	}

	log, ok := r.Context().Value(loggerContextKey{}).(*slog.Logger)
	if !ok || log == nil {
		return
	}
	if httpErr != nil && httpErr.Code == "" && httpErr.Cause == nil {
		return
	}

	cause := err
	if httpErr != nil && httpErr.Cause != nil {
		cause = httpErr.Cause
	}

	attrs := []any{
		"method", r.Method,
		"path", r.URL.Path,
		"status", status,
		"detail", detail,
		"error", cause,
	}
	if code != "" {
		attrs = append(attrs, "code", code)
	}
	if requestID != "" {
		attrs = append(attrs, "request_id", requestID)
	}

	log.Error("http problem", attrs...)
}

func defaultServerErrorCode(status int) string {
	switch status {
	case http.StatusServiceUnavailable:
		return "service_unavailable"
	case http.StatusGatewayTimeout:
		return "gateway_timeout"
	default:
		return "internal_server_error"
	}
}
