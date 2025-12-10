package jwt

import (
	"context"
	"fmt"
	"time"

	models "github.com/Nerfi/instaClone/internal/models/authUser"
	"github.com/go-chi/jwtauth/v5"
)

// Genertatetoken generates a JWT token for a given user ID with an expiration time.
// It takes the user ID, a JWTAuth instance, and the expiration time as parameters.
// It returns the generated token as a string and an error if the token generation fails.

func GenerateToken(user *models.User, auth *jwtauth.JWTAuth, expireTime int64) (string, error) {
	claims := map[string]interface{}{
		"user_id": user.ID,
		"exp":     expireTime,
		"email":   user.Email,
		"iat":     time.Now().Unix(),
		"iss":     "instaClone",
	}

	_, tokestring, err := auth.Encode(claims)
	if err != nil {
		fmt.Println(err, "error genereting the tokens")
		return "", err
	}
	return tokestring, nil

}

// getToken generates and returns new access and refresh tokens for a given user ID.
// It stores the refresh token in the database and returns a TokenResponse containing both tokens.

func GetAuthTokens(user *models.User, auth *jwtauth.JWTAuth) (*models.TokenResponse, error) {
	// generate access token, expires in 30 min
	accessToken, err := GenerateToken(user, auth, time.Now().Add(30*time.Minute).Unix())
	if err != nil {
		return nil, err
	}

	// refresh token, expires in 7 days
	refreshToken, err := GenerateToken(user, auth, time.Now().Add(7*24*time.Hour).Unix())
	if err != nil {
		return nil, err
	}

	return &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func ValidateToken(tokenStr string, auth *jwtauth.JWTAuth) (int, error) {
	// decode el token
	token, err := auth.Decode(tokenStr)
	if err != nil {
		return 0, fmt.Errorf("token inválido: %w", err)
	}

	// verificar la expiracion
	if token.Expiration().Before(time.Now()) {
		return 0, fmt.Errorf("token expirado")
	}

	// extraer el user_id de los claims
	claims, err := token.AsMap(context.Background())
	if err != nil {
		return 0, fmt.Errorf("error extraendo los claims: %w", err)
	}

	// convertir el user_id a int (puede venir como float)
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("error convertiendo el user_id a int")
	}

	// de momento solo extraemos el user_id de los claims , pero deberemos sacar mas cosas para futuro

	return int(userID), nil

}
