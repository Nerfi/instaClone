package services

import (
	"context"
	"fmt"

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

	fmt.Println(posts, "los posts del servicio", err, "<-- error")
	if err != nil {
		return nil, err
	}

	return posts, nil

}
