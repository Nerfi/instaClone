package posts

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	ResModels "github.com/Nerfi/instaClone/internal/models"
	models "github.com/Nerfi/instaClone/internal/models/posts"
	repoSrv "github.com/Nerfi/instaClone/internal/services/posts"
	validator "github.com/Nerfi/instaClone/pkg/validator"
)

type PostsHanlder struct {
	postService repoSrv.PostsService
}

func NewPostsHanlders(service repoSrv.PostsService) *PostsHanlder {
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

func (h *PostsHanlder) GetPost(w http.ResponseWriter, r *http.Request) {
	//extrac the id from the request
	id := r.PathValue("id")

	//convert into int
	idCnvt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//llamamos al servicio
	post, err := h.postService.GetPostByID(r.Context(), idCnvt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//enviamos post si todo esta ok
	ResModels.ResponseWithJSON(w, http.StatusOK, post)

}

func (h *PostsHanlder) CreatePost(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	var postBody models.PostsReqBody
	// adding decoder for unknow fields
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&postBody); err != nil {
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

func (h *PostsHanlder) UpdatePost(w http.ResponseWriter, r *http.Request) {
	// extract the id of the selected post to update
	id := r.PathValue("id")
	idCnvt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	// 2. Parsear el body con los nuevos datos
	var updateData models.PostsReqBody
	// 3. validar los datos
	dec := json.NewDecoder(r.Body)
	// esto lo hacemos por si nos intentan enviar field en la request que no tenemos definidas no entren y salte error
	dec.DisallowUnknownFields()

	if err := dec.Decode(&updateData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate user input data
	// todo rename this function - ValidateReqAuthBody
	if err := validator.ValidateReqAuthBody(updateData); err != nil {
		ResModels.ResponseWithJSON(w, http.StatusBadRequest, err)
		return
	}

	// 4 llamamos al servicio pasando el ID y los datos que queremos actualizar, los que vienen en la request

	post, err := h.postService.UpdatePost(r.Context(), idCnvt, &updateData)

	if err != nil {
		fmt.Println(err.Error(), "el error a la llamada del servicio con los datos que queremos actualizar y el id del elemento")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// devolvemos el post actualizado en la respuesta
	ResModels.ResponseWithJSON(w, http.StatusOK, post)
}
