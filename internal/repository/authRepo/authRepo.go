package repository

import (
	"context"
	"database/sql"
	"fmt"

	models "github.com/Nerfi/instaClone/internal/models/authUser"
)

const (
	INSERT_USER         = "INSERT INTO users(email, password) VALUES(?, ?) "
	COUNT_USER_BY_EMAIL = "SELECT COUNT(EMAIL) FROM users WHERE EMAIL = ?"
	FIND_USER           = "SELECT id, email , password, created_at FROM users WHERE email = ?"
)

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
