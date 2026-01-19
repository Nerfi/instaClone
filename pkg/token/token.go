package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

func GenerateResetToken() (string, string, error) {
	// token real (se envía por email)
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", "", err
	}
	// aqui creamos un token normal y corriente, sin hasing
	token := base64.URLEncoding.EncodeToString(b)

	// hash del token (se guarda en BD)
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	return token, tokenHash, nil
}

// esta funcion recibira el token enviado por el FE(token sin hash) para hashearlo de vuelta y con ese valor buscar en la bbdd si tenemos ese hash, para cambiar la contraseña

func HashToken(token string) string {
	// el token viene en base 64, lo hasheamos igual que en GenerateResetToken
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
