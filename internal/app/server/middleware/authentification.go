package middleware

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"net/http"
	"shortener/internal/app/handler"
	"shortener/internal/app/logging"
	"shortener/internal/app/service/authentification"
)

// AuthMiddleware из куки в контекст
func AuthMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//есть ли кука?
			tokenCookie, err := r.Cookie("Authorization")

			//Если есть ошибка и это не отсутствие куки - прекращаем работу
			if err != nil && !errors.Is(err, http.ErrNoCookie) {
				logging.Sugar.Error(err)
				handler.BadRequest(err, http.StatusInternalServerError)

				return
			}

			//Если куки нет - продолжаем работу
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			//если есть - получить claims из куки
			claims, err := authentification.GetClaims(tokenCookie.Value)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			//Если пользователя нет - продолжаем работу
			if claims.UserID == uuid.Nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), authentification.ContextUserID, claims.UserID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
