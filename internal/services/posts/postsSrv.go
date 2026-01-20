package services

import (
	"context"
	"errors"
	"fmt"

	middleware "github.com/Nerfi/instaClone/internal/handlers/middlewares"
	models "github.com/Nerfi/instaClone/internal/models/posts"
	repo "github.com/Nerfi/instaClone/internal/repository/posts"
)

type PostsService interface {
	// add ctx a este metodo
	GetPosts() ([]*models.Posts, error)
	CreatePost(ctx context.Context, post *models.PostsReqBody) (*models.Posts, error)
	DeletePost(ctx context.Context, id int) error
	GetPostByID(ctx context.Context, id int) (*models.Posts, error)
	UpdatePost(ctx context.Context, id int, updateData *models.PostsReqBody) (*models.Posts, error)
}

// no usar punteros a interfaces, como regla simple
type PostsSrv struct {
	postsrepo repo.PostsRepository
}

func NewPostsSrv(repo repo.PostsRepository) *PostsSrv {
	return &PostsSrv{postsrepo: repo}
}

func (svc *PostsSrv) GetPosts() ([]*models.Posts, error) {
	posts, err := svc.postsrepo.GetPosts(context.Background())

	if err != nil {
		return nil, err
	}

	return posts, nil

}

func (svc *PostsSrv) GetPostByID(ctx context.Context, id int) (*models.Posts, error) {
	post, err := svc.postsrepo.GetPostByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return post, nil
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

func (svc *PostsSrv) UpdatePost(ctx context.Context, id int, updateData *models.PostsReqBody) (*models.Posts, error) {
	// 1 verificar que existe el post que queremos
	pst, err := svc.postsrepo.GetPostByID(ctx, id)
	usrIDCtx := ctx.Value("user_id").(int)
	if err != nil {
		fmt.Println("Error getting post, or does not exist")
		return nil, err
	}

	// verificar que el usuario este autorizado para ejecer esta accion
	usrId := pst.USER_ID
	if usrId != usrIDCtx {
		return nil, errors.New("Forbidden, user can not perform this operation")
	}

	// 2 actualizar el post de la base de datos
	rowsAffected, err := svc.postsrepo.UpdatePost(ctx, id, updateData)

	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("pot not found")
	}

	// 3 obtener el post actualizado para devolverlo
	updatePost, err := svc.postsrepo.GetPostByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return updatePost, nil

}
