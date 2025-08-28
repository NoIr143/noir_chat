package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/noir143/noir_chat/src/shared/exceptions"
)

var Validate = validator.New()

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, err any) {

	switch v := err.(type) {
	case exceptions.InternalException:
		log.Println(v)
		WriteJSON(w, http.StatusInternalServerError, exceptions.HandleExeption("NC_500", http.StatusInternalServerError, "internal_exception"))
	case exceptions.BadRequestException:
		log.Println(v)
		WriteJSON(w, http.StatusBadRequest, exceptions.HandleExeption("NC_400", http.StatusBadRequest, v.Message))
	case exceptions.NotFoundException:
		WriteJSON(w, http.StatusNotFound, exceptions.HandleExeption("NC_404", http.StatusNotFound, "not_found"))
	case exceptions.UnauthorizedException:
		WriteJSON(w, http.StatusUnauthorized, exceptions.HandleExeption("NC_401", http.StatusUnauthorized, "unauthorized"))
	case exceptions.ForbiddenException:
		WriteJSON(w, http.StatusForbidden, exceptions.HandleExeption("NC_403", http.StatusForbidden, "forbidden"))
	case exceptions.InvalidParameterException:
		var invalidParameterResponses []exceptions.InvalidParamterResponse

		for _, e := range v.ValidationErrors {
			invalidParameterResponses = append(invalidParameterResponses, exceptions.InvalidParamterResponse{
				Property: e.Field(),
				Message:  e.Error(),
			})
		}

		WriteJSON(w, http.StatusBadRequest, exceptions.HandleInvalidParameterException("NC_800", http.StatusBadRequest, "invalid_parameter", invalidParameterResponses))
	}
}

func ParseJSON(r *http.Request, v any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(r.Body).Decode(v)
}

func GetTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	tokenQuery := r.URL.Query().Get("token")

	if tokenAuth != "" {
		return tokenAuth
	}

	if tokenQuery != "" {
		return tokenQuery
	}

	return ""
}
