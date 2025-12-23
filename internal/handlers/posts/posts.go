package posts

import (
	"encoding/json"
	ResModels "github.com/Nerfi/instaClone/internal/models"
	models "github.com/Nerfi/instaClone/internal/models/posts"
	repoSrv "github.com/Nerfi/instaClone/internal/services/posts"
	validator "github.com/Nerfi/instaClone/pkg/validator"
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

// TODO add pagination
func (h *PostsHanlder) GetPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.postService.GetPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// success
	ResModels.ResponseWithJSON(w, http.StatusOK, posts)
}

func (h *PostsHanlder) PostPost(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	var postBody models.PostsReqBody

	if err := json.NewDecoder(r.Body).Decode(&postBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate user input data
	// todo rename this function - ValidateReqAuthBody
	if err := validator.ValidateReqAuthBody(postBody); err != nil {
		ResModels.ResponseWithJSON(w, http.StatusBadRequest, err)
		return
	}

	// llamar al servicio si todo o

	post, err := h.postService.CreatePost(r.Context(), &postBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ResModels.ResponseWithJSON(w, http.StatusCreated, post)

}
