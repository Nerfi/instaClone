package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	resModel "github.com/Nerfi/instaClone/internal/models"
	models "github.com/Nerfi/instaClone/internal/models/authUser"
)

const (
	INSERT_USER                = "INSERT INTO users(email, password) VALUES(?, ?) "
	COUNT_USER_BY_EMAIL        = "SELECT COUNT(EMAIL) FROM users WHERE EMAIL = ?"
	FIND_USER                  = "SELECT id, email , password, created_at FROM users WHERE email = ?"
	SELECT_USER                = "SELECT * FROM users WHERE id = ?"
	INSERT_REFRESH_TOKEN       = "INSERT INTO refresh_tokens_table( user_id, token, expires_at) VALUES(?, ?, ?)"
	DELETE_TOKEN_REFRESH_TABLE = "DELETE FROM refresh_tokens_table WHERE user_id = ?"
	FIND_REFRESH_TOKEN         = "SELECT user_id, expires_at FROM refresh_tokens_table WHERE token = ?"
	DELETE_SPECIFIC_TOKEN      = "DELETE FROM refresh_tokens_table WHERE token = ?"
	FIND_USER_BY_EMAIL         = "SELECT id, email FROM users WHERE EMail = ?"
)

//TODO: create interface para repo y servicio

type AuthRepo struct {
	db *sql.DB
}

func NewAuthRepo(db *sql.DB) *AuthRepo {
	return &AuthRepo{db: db} // bring this from config
}

func (r *AuthRepo) CreateUser(ctx context.Context, user *models.User) (int64, error) {

	// aqui va la conexion con la bbdd junto con la query para poder insertar el usuario
	var row int
	err := r.db.QueryRow(COUNT_USER_BY_EMAIL, user.Email).Scan(&row)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println(err, "error getting user")
		return 0, err
	}
	// si row != 0 es que ya hay un usuario con ese emial
	if row != 0 {
		return 0, fmt.Errorf("email already exists")
	}

	// insertamos el usuario en la bbdd con la password hasheada
	result, err := r.db.Exec(INSERT_USER, user.Email, user.Password)
	if err != nil {
		fmt.Println(err, "error inserting user")
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println(err, "error getting user")
		return 0, err
	}

	// if all went good
	return id, nil
}

func (r *AuthRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	// cuando el struct no esta inicializado como en este caso , no tenemos que usar punteros (*) sino que debemos inicializar el struct en memoria , como abajo
	// recordar tambien que debemos de rellenar el metodo Scan con el numero de columnas selecciondas en la query que ejeuctamos, aqui son dos parametros para Scan
	// ya que nuestra query selecciona id , email, password y createdAt
	var row models.User

	user := r.db.QueryRow(FIND_USER, email)
	if err := user.Scan(&row.ID, &row.Email, &row.Password, &row.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user with such email")
		}
		return nil, err
	}
	return &row, nil
}

func (r *AuthRepo) LogOutUser(ctx context.Context, userId int) (*resModel.Response, error) {
	_, err := r.db.Exec(DELETE_TOKEN_REFRESH_TABLE, userId)
	if err != nil {
		return nil, err
	}
	return &resModel.Response{Success: true, StatusCode: 200, Data: "user logged out"}, nil
}

func (r *AuthRepo) Profile(ctx context.Context, id int) (*models.User, error) {
	var userM models.User
	user := r.db.QueryRow(SELECT_USER, id)

	if err := user.Scan(&userM.ID, &userM.Email, &userM.Password, &userM.CreatedAt); err != nil {
		return nil, err
	}

	return &userM, nil
}

// insertamos el refresh token en la bbdd
func (r *AuthRepo) InsertRefreshToken(ctx context.Context, token string, userID int, expiresAt time.Time) error {
	_, err := r.db.Exec(INSERT_REFRESH_TOKEN, userID, token, expiresAt)
	return err
}

// refresh token logic
func (r *AuthRepo) GetRefreshToken(ctx context.Context, token string) (int, time.Time, error) {
	var userID int
	var expiresAt time.Time
	err := r.db.QueryRowContext(ctx, FIND_REFRESH_TOKEN, token).Scan(&userID, &expiresAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, time.Time{}, fmt.Errorf("token not found")
		}
		return 0, time.Time{}, err
	}

	return userID, expiresAt, nil
}

func (r *AuthRepo) DeleteRefreshToken(ctx context.Context, token string) error {
	_, err := r.db.Exec(DELETE_SPECIFIC_TOKEN, token)
	return err
}

func (r *AuthRepo) GetUserById(ctx context.Context, userID int) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(SELECT_USER, userID).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found")
		}
		return nil, err
	}
	return &user, nil

}

func (r *AuthRepo) FindUserByEmail(ctx context.Context, email string) (*models.ChangePasswordUser, error) {
	var usrcPssChng models.ChangePasswordUser
	err := r.db.QueryRow(FIND_USER_BY_EMAIL, email).Scan(&usrcPssChng.ID, &usrcPssChng.Email)
	if err != nil {
		return nil, err
	}
	return &usrcPssChng, nil
}
