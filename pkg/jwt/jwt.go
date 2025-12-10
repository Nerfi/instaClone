package jwt

import (
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
