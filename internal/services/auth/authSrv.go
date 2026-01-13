package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	resModel "github.com/Nerfi/instaClone/internal/models"
	"github.com/Nerfi/instaClone/internal/models/authUser"
	authRepo "github.com/Nerfi/instaClone/internal/repository/authRepo"
	jwt "github.com/Nerfi/instaClone/pkg/jwt"
	"github.com/Nerfi/instaClone/pkg/security"
	validator "github.com/Nerfi/instaClone/pkg/validator"
	"github.com/go-chi/jwtauth/v5"
)

type AuthSrv struct {
	authrepo *authRepo.AuthRepo
	auth     *jwtauth.JWTAuth
}

//TODO: create interface for service and repo

func NewAuthSrv(repo *authRepo.AuthRepo, auth *jwtauth.JWTAuth) *AuthSrv {
	return &AuthSrv{authrepo: repo, auth: auth}
}

func (svc *AuthSrv) CreteUser(ctx context.Context, body *models.AuthReqBody) (*models.User, error) {
	// validamos el input antes de hacer nada o continuar
	if msg := validator.ValidateReqAuthBody(*body); msg != nil {
		return nil, fmt.Errorf(strings.Join(msg, ", "))
	}

	// hash de la contraseña antes de guardar el usuario tambien comprobamos la longitud y otras caracteristicas de la password
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

	// validamos los datos enviados por el usuario
	if msg := validator.ValidateReqAuthBody(*body); msg != nil {
		return nil, fmt.Errorf(strings.Join(msg, ", "))
	}
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

	// guardamos el refresh token en la base de datos
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	if err = svc.authrepo.InsertRefreshToken(ctx, tokens.RefreshToken, dbUser.ID, expiresAt); err != nil {
		return nil, err
	}

	return tokens, nil

}

func (svc *AuthSrv) LogOutUser(ctx context.Context, userID int) (*resModel.Response, error) {
	usrRes, err := svc.authrepo.LogOutUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return usrRes, nil

}

func (svc *AuthSrv) Profile(ctx context.Context, userId int) (*models.User, error) {
	// buscamos al usuario en la bbdd en base al id pasado
	usr, err := svc.authrepo.Profile(ctx, userId)
	if err != nil {
		return nil, err
	}
	return usr, nil
}

func (svc *AuthSrv) CheckRefreshTokenValid(ctx context.Context, token string) (*models.TokenResponse, error) {

	userId, expiresAt, err := svc.authrepo.GetRefreshToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if time.Now().After(expiresAt) {
		return nil, fmt.Errorf("refresh token expired")
	}

	// buscamos al usuario para llamar al metodo que crea los tokens y los devuelve
	fullUser, err := svc.authrepo.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	// generar los nuevos tokens de acceso y refresco
	tokens, err := jwt.GetAuthTokens(fullUser, svc.auth)
	if err != nil {
		return nil, err
	}

	// generar nuevos tokens y devolverlos al handler
	return tokens, nil
}

func (svc *AuthSrv) FindUserByEmail(ctx context.Context, email string) (*models.ChangePasswordUser, error) {
	crrntUsr, err := svc.authrepo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return crrntUsr, nil
}
