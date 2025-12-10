package auth

import (
	"encoding/json"
	"net/http"

	"github.com/Nerfi/instaClone/internal/models/authUser"
	authsrvc "github.com/Nerfi/instaClone/internal/services/auth"
	"github.com/gorilla/csrf"
)

type AuthHanlders struct {
	authservice *authsrvc.AuthSrv
}

func NewAuthHanlders(service *authsrvc.AuthSrv) *AuthHanlders {
	return &AuthHanlders{
		authservice: service,
	}
}

func (h *AuthHanlders) CreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var bodyReq *models.AuthReqBody

	if err := json.NewDecoder(r.Body).Decode(&bodyReq); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("please provied valid input ")
		return
	}

	// llamar al servicio
	user, err := h.authservice.CreteUser(r.Context(), bodyReq)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// devolviendo en caso de que todo haya ido bien
	w.Header().Set("X-CSRF-Token", csrf.Token(r))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)

}

func (h *AuthHanlders) LoginUser(w http.ResponseWriter, r *http.Request) {
	var bodyReq *models.AuthReqBody
	if err := json.NewDecoder(r.Body).Decode(&bodyReq); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("please provied valid input ")
		return
	}

	tokens, err := h.authservice.LoginUser(r.Context(), bodyReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// guardamos los tokens en las cookieshttponly
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   604800, // 7 días
	})
	// seteamos el access token en las cookies tambien
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokens.AccessToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   900, // 15 minutos
	})

	// Devolver access token en el body
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful",
	})

}
