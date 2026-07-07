package core_http_middleware

import (
	"fmt"
	"net/http"
	"strings"

	core_auth "github.com/alisupurov/todoApp-golang/internal/core/auth"
	core_errors "github.com/alisupurov/todoApp-golang/internal/core/errors"
	core_logger "github.com/alisupurov/todoApp-golang/internal/core/logger"
	core_http_context "github.com/alisupurov/todoApp-golang/internal/core/transport/http/context"
	core_http_response "github.com/alisupurov/todoApp-golang/internal/core/transport/http/response"
)

func JWT(config core_auth.Config) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				responseHandler.ErrorResponse(
					fmt.Errorf("missing authorization header: %w", core_errors.ErrUnauthorized),
					"unauthorized",
				)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := core_auth.ValidateToken(tokenString, config)
			if err != nil {
				responseHandler.ErrorResponse(err, "unauthorized")
				return
			}

			ctx = core_http_context.WithAccountID(ctx, claims.AccountID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
