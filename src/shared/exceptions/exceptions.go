package exceptions

import (
	_ "embed"
	"encoding/json"

	"github.com/go-playground/validator/v10"
	"github.com/noir143/noir_chat/src/configs"
)

//go:embed messages.json
var messagesData []byte

//go:embed internal_messages.json
var internalMessagesData []byte

type InternalException struct {
	Error error
}

type BadRequestException struct {
	ErrorId string
	Message string
	Error   error
}

type InvalidParameterException struct {
	ValidationErrors validator.ValidationErrors
}

type NotFoundException struct {
	Error error
}

type UnauthorizedException struct {
	Error error
}

type ForbiddenException struct {
	Error error
}

type InvalidParamterResponse struct {
	Property string `json:"property"`
	Message  string `json:"message"`
}

type ErrorResponse struct {
	ErrorId string                    `json:"error_id"`
	Status  int                       `json:"status"`
	Message string                    `json:"message"`
	Errors  []InvalidParamterResponse `json:"errors"`
}

func HandleExeption(errorId string, status int, err string) ErrorResponse {
	return ErrorResponse{
		ErrorId: errorId,
		Status:  status,
		Message: getMesage(err),
	}
}

func HandleInvalidParameterException(errorId string, status int, err string, errors []InvalidParamterResponse) ErrorResponse {
	return ErrorResponse{
		ErrorId: errorId,
		Status:  status,
		Message: getMesage(err),
		Errors:  errors,
	}
}

func getMesage(key string) string {
	var data map[string]interface{}

	src := messagesData
	if configs.EnvConfigs.ENABLE_DEBUB_MODE {
		src = internalMessagesData
	}

	if err := json.Unmarshal(src, &data); err != nil {
		panic(err)
	}

	if val, ok := data[key].(string); ok {
		return val
	}
	return ""
}
