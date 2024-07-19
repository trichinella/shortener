package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"shortener/internal/app/handler"
	"shortener/internal/app/logging"
	"shortener/internal/app/service/authentification"
	"time"
)

func AuthMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//есть ли кука?
			tokenCookie, err := r.Cookie("token")

			//Если есть ошибка и это не отсутствие куки - прекращаем работу
			if err != nil && !errors.Is(err, http.ErrNoCookie) {
				logging.Sugar.Error(err)
				handler.BadRequest(err, http.StatusInternalServerError)

				return
			}

			//Если куки нет - получить токен
			if err != nil {
				createNewToken(next, w, r)
				return
			}

			//если есть - получить claims из куки
			claims, err := authentification.GetClaims(tokenCookie.Value)
			if err != nil {
				createNewToken(next, w, r)
				return
			}

			if claims.UserID == uuid.Nil {
				handler.BadRequest(fmt.Errorf("unauthorized"), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "UserID", claims.UserID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func createNewToken(next http.Handler, w http.ResponseWriter, r *http.Request) {
	signedToken, err := authentification.BuildJWTString()

	if err != nil {
		logging.Sugar.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   signedToken,
		Expires: time.Now().Add(3 * time.Hour),
	})

	claims, err := authentification.GetClaims(signedToken)
	if err != nil {
		handler.BadRequest(err, http.StatusBadRequest)
		return
	}

	ctx := context.WithValue(r.Context(), "UserID", claims.UserID)
	r = r.WithContext(ctx)

	next.ServeHTTP(w, r)
}
