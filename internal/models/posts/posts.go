package models

import "time"

type Posts struct {
	ID        int       `json:"id"`
	USER_ID   int       `json:"user_id"`
	Caption   string    `json:"caption"`
	Image_url string    `json:"image_url"` // maybe not this or update in order to have more photos, we ll see
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostsReqBody struct {
	Caption   string `json:"caption" validate:"required,min=1,max=140"`
	Image_url string `json:"image_url" validate:"required,url"` // maybe url as validate is not the best, check it later
}
