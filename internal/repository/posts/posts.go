package repository

import (
	"context"
	"database/sql"
	"fmt"

	models "github.com/Nerfi/instaClone/internal/models/posts"
)

const (
	SELECT_POSTS        = "SELECT * FROM posts"
	CREATE_POST         = "INSERT INTO posts(user_id, caption, image_url) VALUES(?, ?, ?)"
	DELETE_POST_BY_ID   = "DELETE FROM posts WHERE id = ?"
	CHECK_OWNER_OS_POST = "SELECT user_id FROM posts WHERE id = ?"
	GET_SINGLE_POST     = "SELECT * FROM posts WHERE id = ?"
	UPDATE_POST         = "UPDATE posts SET caption = ?, image_url = ?, updated_at = NOW() WHERE id = ?"
)

type PostsRepository interface {
	GetPosts(ctx context.Context) ([]*models.Posts, error)
	GetPostByID(ctx context.Context, id int) (*models.Posts, error)
	CreatePost(ctx context.Context, post *models.PostsReqBody) (*models.Posts, error)
	DeletePost(ctx context.Context, id int) error
	GetPostOwner(ctx context.Context, id int) (int, error)
	UpdatePost(ctx context.Context, id int, post *models.PostsReqBody) (int64, error)
}

type PostsRepo struct {
	db *sql.DB
}

func NewPostsRepo(db *sql.DB) *PostsRepo {
	return &PostsRepo{db: db}
}

func (r *PostsRepo) GetPosts(ctx context.Context) ([]*models.Posts, error) {
	rows, err := r.db.Query(SELECT_POSTS)
	if err != nil {
		return nil, fmt.Errorf("error getting posts: %w", err)
	}

	defer rows.Close() // important , always close rows

	var posts []*models.Posts

	for rows.Next() {
		var post models.Posts
		if err := rows.Scan(&post.ID, &post.USER_ID, &post.Caption, &post.Image_url, &post.CreatedAt, &post.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	// important check for errors during iteration

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil

}

func (r *PostsRepo) GetPostByID(ctx context.Context, id int) (*models.Posts, error) {
	var p models.Posts
	err := r.db.QueryRow(GET_SINGLE_POST, id).Scan(&p.ID, &p.USER_ID, &p.Caption, &p.Image_url, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("error getting single post: %w", err)
	}
	return &p, nil

}

func (r *PostsRepo) CreatePost(ctx context.Context, post *models.PostsReqBody) (*models.Posts, error) {
	result, err := r.db.Exec(CREATE_POST, ctx.Value("user_id"), post.Caption, post.Image_url)

	fmt.Printf("USER ID IN CONTEXT: %#v\n", ctx.Value("user_id"))

	if err != nil {
		return nil, fmt.Errorf("db error creating post %w", err)
	}
	returnedPost := &models.Posts{
		USER_ID:   ctx.Value("user_id").(int),
		Caption:   post.Caption,
		Image_url: post.Image_url,
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	returnedPost.ID = int(lastID)
	return returnedPost, nil
}

func (r *PostsRepo) DeletePost(ctx context.Context, id int) error {
	result, err := r.db.Exec(DELETE_POST_BY_ID, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no post found with id %d", id)
	}

	return nil
}

func (r *PostsRepo) GetPostOwner(ctx context.Context, id int) (int, error) {
	var ownerId int
	err := r.db.QueryRow(CHECK_OWNER_OS_POST, id).Scan(&ownerId)
	if err != nil {
		return 0, err
	}
	return ownerId, nil
}

func (r *PostsRepo) UpdatePost(ctx context.Context, id int, post *models.PostsReqBody) (int64, error) {
	result, err := r.db.Exec(UPDATE_POST, post.Caption, post.Image_url, id)
	if err != nil {
		return 0, err
	}

	// verificamos cuantas filas se actualizarion
	// RowsAffected returns the number of rows affected by an update, insert, or delete.
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, nil
}
