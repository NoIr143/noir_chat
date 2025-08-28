package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/noir143/noir_chat/src/configs"
	"github.com/noir143/noir_chat/src/shared/exceptions"
	"github.com/noir143/noir_chat/src/shared/utils"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := utils.GetTokenFromRequest(r)
		if tokenString == "" {
			utils.WriteError(w, exceptions.UnauthorizedException{Error: fmt.Errorf("missing token")})
			return
		}

		if strings.HasPrefix(strings.ToLower(tokenString), "bearer ") {
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
			tokenString = strings.TrimPrefix(tokenString, "bearer ")
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(configs.EnvConfigs.JWT_SECRET), nil
		})

		if err != nil {
			utils.WriteError(w, exceptions.UnauthorizedException{Error: err})
			return
		}

		if !token.Valid {
			utils.WriteError(w, exceptions.UnauthorizedException{Error: fmt.Errorf("invalid token")})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.WriteError(w, exceptions.UnauthorizedException{Error: fmt.Errorf("invalid claims")})
			return
		}

		if expRaw, ok := claims["expiresAt"]; ok {
			var expUnix int64
			switch v := expRaw.(type) {
			case float64:
				expUnix = int64(v)
			case int64:
				expUnix = v
			case string:
				if n, e := strconv.ParseInt(v, 10, 64); e == nil {
					expUnix = n
				}
			}
			if expUnix == 0 || time.Now().Unix() > expUnix {
				utils.WriteError(w, exceptions.UnauthorizedException{Error: fmt.Errorf("token expired")})
				return
			}
		}

		var userIDStr string
		if v, ok := claims["userID"].(string); ok {
			userIDStr = v
		} else if vFloat, ok := claims["userID"].(float64); ok {
			userIDStr = strconv.Itoa(int(vFloat))
		}
		if userIDStr == "" {
			utils.WriteError(w, exceptions.UnauthorizedException{Error: fmt.Errorf("missing user id")})
			return
		}

		userIDInt, convErr := strconv.Atoi(userIDStr)
		if convErr != nil {
			utils.WriteError(w, exceptions.UnauthorizedException{Error: fmt.Errorf("invalid user id")})
			return
		}

		ctx := context.WithValue(r.Context(), utils.UserKey, userIDInt)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
