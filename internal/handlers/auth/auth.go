package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	usrMiddle "github.com/Nerfi/instaClone/internal/handlers/middlewares"
	ResModels "github.com/Nerfi/instaClone/internal/models"
	"github.com/Nerfi/instaClone/internal/models/authUser"
	authsrvc "github.com/Nerfi/instaClone/internal/services/auth"
	"github.com/Nerfi/instaClone/pkg/security"
	"github.com/Nerfi/instaClone/pkg/token"
	validation "github.com/Nerfi/instaClone/pkg/validator"
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
		ResModels.ResponseWithJSON(w, http.StatusBadRequest, "please provide valid input ")
		return
	}

	// validando el req body
	if err := validation.ValidateReqAuthBody(*bodyReq); err != nil {
		ResModels.ResponseWithJSON(w, http.StatusBadRequest, err)
		return
	}

	// llamar al servicio
	user, err := h.authservice.CreteUser(r.Context(), bodyReq)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// devolviendo en caso de que todo haya ido bien

	// TODO implemente csrf token
	w.Header().Set("X-CSRF-Token", csrf.Token(r))
	ResModels.ResponseWithJSON(w, http.StatusCreated, user)

}

func (h *AuthHanlders) LoginUser(w http.ResponseWriter, r *http.Request) {
	var bodyReq models.AuthReqBody
	if err := json.NewDecoder(r.Body).Decode(&bodyReq); err != nil {
		ResModels.ResponseWithJSON(w, http.StatusBadRequest, "please provide valid input ")

		return
	}

	// validation user input
	if err := validation.ValidateReqAuthBody(bodyReq); err != nil {
		ResModels.ResponseWithJSON(w, http.StatusBadRequest, err)
		return
	}

	tokens, err := h.authservice.LoginUser(r.Context(), &bodyReq)
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
		Path:     "/auth/refresh",
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
	ResModels.ResponseWithJSON(w, http.StatusOK, "Login successful")

}

func (h *AuthHanlders) LogoutUser(w http.ResponseWriter, r *http.Request) {
	// extraer el id del usuario del context
	userId, ok := usrMiddle.GetUserIdFromContext(r.Context())
	if !ok {
		http.Error(w, "no user found", http.StatusUnauthorized)
		return
	}

	resut, err := h.authservice.LogOutUser(r.Context(), userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// limpiamos las cookies de access_token y refresh_token
	// https://stackoverflow.com/questions/27671061/how-to-delete-cookie

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/auth/refresh",
		MaxAge:   -1,
	})

	ResModels.ResponseWithJSON(w, http.StatusOK, resut)
}

func (h *AuthHanlders) Profile(w http.ResponseWriter, r *http.Request) {
	userId, ok := usrMiddle.GetUserIdFromContext(r.Context())
	if !ok {
		http.Error(w, "no user found", http.StatusUnauthorized)
		return
	}
	// hablamos con el servicio para extraer los datos del usuario
	usrPfl, err := h.authservice.Profile(r.Context(), userId)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	ResModels.ResponseWithJSON(w, http.StatusOK, usrPfl)

}

// RefreshToken handles the token refresh process for authenticated users.
// It retrieves the JWT cookie from the request, validates it, and generates new access and refresh tokens.
// If the JWT cookie is missing or invalid, it responds with an unauthorized status.
// If the token generation is successful, it sets the new refresh token in the cookie and responds with the new access token.

// @param w http.ResponseWriter - the response writer to send the response
// @param r *http.Request - the incoming HTTP request containing the JWT cookie

func (h *AuthHanlders) RefreshToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil || cookie.Value == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	tokenValue := cookie.Value
	newUsrTokens, err := h.authservice.CheckRefreshTokenValid(r.Context(), tokenValue)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// seteamos de nuevo los tokens en las cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    newUsrTokens.AccessToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   900, // 15 minutos
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newUsrTokens.RefreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/auth/refresh",
		MaxAge:   604800, // 7 días
	})

	// response si todo ha ido bien
	ResModels.ResponseWithJSON(w, http.StatusOK, "tokens refreshed successfully")

}

// forgot password function
// debemos recibir el email enviado y buscar si ese usuario existe en nuestro sistema
// debemos generar token si el usuario existe
// guardar el token y enviar el email con el link(token no hasheado)
func (h *AuthHanlders) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	// read the email sent and check if user exist in DB
	usrMail := r.FormValue("email")
	// Respuesta anti‑enumeración, lo hacemos para que si un atacante intenta atacarnos no sepa si
	// el emial existe o no
	genericResponse := map[string]string{
		"message": "If the email exists, a reset link will be sent to you",
	}
	// buscamos el usuario para ver si existe en bbdd( si no existe seguimos igual)
	usrc, err := h.authservice.FindUserByEmail(r.Context(), usrMail)
	if err != nil || usrc == nil {
		ResModels.ResponseWithJSON(w, http.StatusOK, genericResponse)
		return
	}

	// generar token seguro
	tokenSecure, tokenHash, err := token.GenerateResetToken()
	if err != nil {
		ResModels.ResponseWithJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	// guardar token en la bbdd, para despues buscarlo y saber si existe o no
	expiresAt := time.Now().Add(30 * time.Minute)
	err = h.authservice.SavePasswordResetToken(r.Context(), usrc.ID, tokenHash, expiresAt)
	if err != nil {
		ResModels.ResponseWithJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	//construir URL, el tokenSecure aqui es el token raw, sin hashing, este token es el que tendremos que ,mediante una funcion, hashear para encontrar el valor en bbdd del token
	// resetURL := fmt.Sprintf("%s/reset-password?token=%s", h.config.FrontendURL, tokenSecure)
	fmt.Sprintf("LLEGAMOS A FORGOT PASSWORD URL")
	fmt.Sprintf("%s/reset-password?token=%s", "loclahost:3000/reset", tokenSecure)
	fmt.Println(tokenSecure, "%s/reset-password?token=%s", "loclahost:3000/reset", tokenSecure)

	// enviar email( aunque falle no se revela nada)
	// TODO: implement email service and logic
	// _ = h.emailService.SendPasswordResetEmail(usrc.Email, resetURL)

	// respuesta generica
	ResModels.ResponseWithJSON(w, http.StatusOK, genericResponse)

}
func (h *AuthHanlders) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req models.ResetPasswordRequest
	// parseamos los datos
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ResModels.ResponseWithJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
		return
	}

	// validar que los campos no esten vacios
	// validation user input
	if err := validation.ValidateReqAuthBody(req); err != nil {
		ResModels.ResponseWithJSON(w, http.StatusBadRequest, err)
		return
	}

	// convertir el token recibido a su hash , para buscarlo en la bbdd
	tokenHash := token.HashToken(req.Token)

	// validar el token y obtener el userID
	userID, err := h.authservice.ValidateResetToken(r.Context(), tokenHash)
	if err != nil {
		fmt.Println(err.Error())
		ResModels.ResponseWithJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid or expired token 1",
		})
		return
	}

	// hashear la nueva contraseña para guardarla en la base de datos
	hashedNewUserPassword, err := security.HashPassword(req.NewPassword)
	if err != nil {
		// No revelamos detalles específicos del error (seguridad)
		ResModels.ResponseWithJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	// actualizar la contraseña del usuario en su tabla
	err = h.authservice.UpdateUserPassword(r.Context(), userID, hashedNewUserPassword)
	if err != nil {
		ResModels.ResponseWithJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to updated password",
		})
		return
	}

	// marcar el token como usado
	err = h.authservice.MarkTokenAsUsed(r.Context(), tokenHash)
	if err != nil {
		// logueamos el error pero no paramos la ejecucion , ya que si hemos llegado aqui todo ha ido bien, este paso es solo limpieza en la bbdd
		log.Printf("Warning: failed to mark token as used: %v", err)
	}

	// respuesta existosa
	ResModels.ResponseWithJSON(w, http.StatusOK, map[string]string{
		"message": "Password has been reset successfully",
	})
}
