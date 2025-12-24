package posts

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	ResModels "github.com/Nerfi/instaClone/internal/models"
	models "github.com/Nerfi/instaClone/internal/models/posts"
	repoSrv "github.com/Nerfi/instaClone/internal/services/posts"
	validator "github.com/Nerfi/instaClone/pkg/validator"
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

	// llamar al servicio si todo ok

	post, err := h.postService.CreatePost(r.Context(), &postBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ResModels.ResponseWithJSON(w, http.StatusCreated, post)

}

func (h *PostsHanlder) DeletePost(w http.ResponseWriter, r *http.Request) {
	// extract the param selected and check if we can do such action
	// only owner of the post should delete it
	idStr := r.PathValue("id")
	//convert into int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = h.postService.DeletePost(r.Context(), id)
	if err != nil {

		if err.Error() == "unauthorized" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if err.Error() == "Forbidden" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "post not found", http.StatusNotFound)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
