package posts

import (
	ResModels "github.com/Nerfi/instaClone/internal/models"
	repoSrv "github.com/Nerfi/instaClone/internal/services/posts"
	"net/http"
)

type PostsHanlder struct {
	postService *repoSrv.PostsSrv
}

func NewPostsHanlders(service *repoSrv.PostsSrv) *PostsHanlder {
	return &PostsHanlder{
		postService: service,
	}
}

func (h *PostsHanlder) GetPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.postService.GetPosts()
	if err != nil {
		ResModels.ResponseWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	// success
	ResModels.ResponseWithJSON(w, http.StatusOK, posts)
}
