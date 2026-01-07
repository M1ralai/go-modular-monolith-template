package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/common/stype"
	"github.com/go-playground/validator/v10"
)

type ctxKey string

const RoleKey ctxKey = "role"
const UsernameKey ctxKey = "username"
const UserIDKey ctxKey = "user_id"

func ReadJson[T any](r *http.Request, validate *validator.Validate) (T, error) {
	var res T
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		return res, err
	}
	if validate != nil {
		err = validate.Struct(res)
		if err != nil {
			return res, err
		}
	}
	return res, nil
}

func WriteJson(w http.ResponseWriter, data interface{}, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := SuccessResponse(data, message)
	json.NewEncoder(w).Encode(resp)
}

func Return(w http.ResponseWriter, statusCode int, resp stype.APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}

func SuccessResponse(data interface{}, message string) stype.APIResponse {
	return stype.APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Error:     nil,
		Timestamp: getCurrentTimestamp(),
	}
}

func ErrorResponse(code, message, details string) stype.APIResponse {
	return stype.APIResponse{
		Success: false,
		Message: message,
		Data:    nil,
		Error: &stype.ErrorDetail{
			Code:    code,
			Details: details,
		},
		Timestamp: getCurrentTimestamp(),
	}
}

func getCurrentTimestamp() string {
	return fmt.Sprintf("%s", time.Now().Format(time.RFC3339))
}

func GetUsernameFromContext(ctx interface{}) string {
	if c, ok := ctx.(interface{ Value(any) any }); ok {
		if username, ok := c.Value(UsernameKey).(string); ok {
			return username
		}
	}
	return "unknown"
}

func GetUserIDFromContext(ctx interface{}) int {
	if c, ok := ctx.(interface{ Value(any) any }); ok {
		if id, ok := c.Value(UserIDKey).(int); ok {
			return id
		}
	}
	return 0
}

func ReturnError(w http.ResponseWriter, code, message, details string) {
	var status int
	switch code {
	case "VALIDATION_ERROR", "BAD_REQUEST":
		status = http.StatusBadRequest
	case "UNAUTHORIZED":
		status = http.StatusUnauthorized
	case "FORBIDDEN":
		status = http.StatusForbidden
	case "NOT_FOUND":
		status = http.StatusNotFound
	default:
		status = http.StatusInternalServerError
	}

	resp := ErrorResponse(code, message, details)
	Return(w, status, resp)
}
