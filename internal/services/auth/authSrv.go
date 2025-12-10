package services

import (
	"context"
	"fmt"
	"time"

	"github.com/Nerfi/instaClone/internal/models/authUser"
	authRepo "github.com/Nerfi/instaClone/internal/repository/authRepo"
	jwt "github.com/Nerfi/instaClone/pkg/jwt"
	"github.com/Nerfi/instaClone/pkg/security"
	"github.com/go-chi/jwtauth/v5"
)

type AuthSrv struct {
	authrepo *authRepo.AuthRepo
	auth     *jwtauth.JWTAuth
}

func NewAuthSrv(repo *authRepo.AuthRepo, auth *jwtauth.JWTAuth) *AuthSrv {
	return &AuthSrv{authrepo: repo, auth: auth}
}

func (svc *AuthSrv) CreteUser(ctx context.Context, body *models.AuthReqBody) (*models.User, error) {
	// hash de la contraseña antes de guardar el usuario
	hashPassword, err := security.HashPassword(body.Password)
	if err != nil {
		return nil, err
	}

	// con la password hasheada rellenamos el struct que contiene los campos con los que vamos a rellenar la struc para crear nuestro usuario en bbdd
	user := &models.User{Email: body.Email, Password: hashPassword, CreatedAt: time.Now()}

	// llamar al repo de create
	id, err := svc.authrepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	// si todo sale ok
	user.ID = int(id)
	return user, nil

}

func (svc *AuthSrv) LoginUser(ctx context.Context, body *models.AuthReqBody) (*models.TokenResponse, error) {
	// 1 buscar usuario en bbdd
	dbUser, err := svc.authrepo.GetUserByEmail(ctx, body.Email)
	if err != nil {
		return nil, err
	}
	// 2 comprobar la contraseña del usuario con el hash
	ok, err := security.ComparePassword(dbUser.Password, body.Password)
	if !ok || err != nil {
		return nil, fmt.Errorf("invalid credentials-->", err)
	}

	// add tokens(refresh , access) and set it into cookies
	tokens, err := jwt.GetAuthTokens(dbUser, svc.auth)
	if err != nil {
		fmt.Println(err, "error generating tokens")
		return nil, err
	}

	return tokens, nil

}
