package services

import (
	"context"

	models "github.com/Nerfi/instaClone/internal/models/posts"
	repo "github.com/Nerfi/instaClone/internal/repository/posts"
)

type PostsSrv struct {
	postsrepo *repo.PostsRepo
}

func NewPostsSrv(repo *repo.PostsRepo) *PostsSrv {
	return &PostsSrv{postsrepo: repo}
}

func (svc *PostsSrv) GetPosts() ([]*models.Posts, error) {
	posts, err := svc.postsrepo.GetPosts(context.Background())

	if err != nil {
		return nil, err
	}

	return posts, nil

}

func (svc *PostsSrv) CreatePost(ctx context.Context, post *models.PostsReqBody) (*models.Posts, error) {
	createdPost, err := svc.postsrepo.CreatePost(ctx, post)
	if err != nil {
		return nil, err
	}

	return createdPost, nil
}
