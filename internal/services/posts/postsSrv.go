package services

import (
	"context"
	"errors"

	middleware "github.com/Nerfi/instaClone/internal/handlers/middlewares"
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

func (svc *PostsSrv) DeletePost(ctx context.Context, id int) error {
	userID, ok := middleware.GetUserIdFromContext(ctx)
	if !ok {
		return errors.New("unauthorized")
	}
	//check ownership, if the post user_id is not the same as the current user this wont work
	ownerId, err := svc.postsrepo.GetPostOwner(ctx, id)

	if err != nil {
		return err
	}
	if userID != ownerId {
		return errors.New("Forbidden")
	}

	return svc.postsrepo.DeletePost(ctx, id)

}
