package middlewares

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Nerfi/instaClone/internal/config"
	security "github.com/Nerfi/instaClone/pkg/jwt"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the token from cookie
		cookie, err := r.Cookie("access_token")
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		tokenStr := cookie.Value

		if tokenStr == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// obtener la config de jwt
		auth := config.NewAppConfig().Auth
		// validar el token

		user, err := security.ValidateToken(tokenStr, auth)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Agregar el userID and email al contexto
		// https://stackoverflow.com/questions/40379960/context-withvalue-how-to-add-several-key-value-pairs
		ctx := context.WithValue(r.Context(), "user_id", user.UserID)
		ctx = context.WithValue(ctx, "user_email", user.Email)

		// Continuar con el siguiente handler
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func OwnerOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// cogemos el id del usuario que esta en el context /logueado
		userId, ok := GetUserIdFromContext(r.Context())
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// extraer el id del path
		requestID := r.PathValue("id")
		requestUserID, err := strconv.Atoi(requestID)
		if err != nil {
			http.Error(w, "error converting id to int", http.StatusBadRequest)
			return
		}

		// check validation para que el usuario acceda a sus propios datos

		if userId != requestUserID {
			http.Error(w, "unauthorized from ids", http.StatusForbidden)
			return
		}

		// todo bien, continuamos
		next.ServeHTTP(w, r)
	})
}

// extraer el ID del contexto
func GetUserIdFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value("user_id").(int)
	return userID, ok
}

// maybe move this to other file
func ChainMiddleware(mw ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(final http.Handler) http.Handler {
		for i := len(mw) - 1; i >= 0; i-- {
			final = mw[i](final)
		}
		return final
	}
}
