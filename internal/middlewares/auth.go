package middlewares

import (
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"os"
)

var (
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("session")
		if err != nil {
			http.Error(w, "unauthorized", 401)
			return
		}

		token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "unauthorized", 401)
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		// сохраняем claims в контекст запроса
		ctx := WithClaims(r.Context(), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
