package repository

import (
	"context"
	"database/sql"
	"fmt"

	models "github.com/Nerfi/instaClone/internal/models/posts"
)

const (
	SELECT_POSTS = "SELECT * FROM posts"
)

type PostsRepo struct {
	db *sql.DB
}

func NewPostsRepo(db *sql.DB) *PostsRepo {
	return &PostsRepo{db: db}
}

func (r *PostsRepo) GetPosts(ctx context.Context) ([]*models.Posts, error) {
	rows, err := r.db.Query(SELECT_POSTS)
	fmt.Println(rows, "los posts que deberian de venirme ")
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
