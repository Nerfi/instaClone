package models

type Posts struct {
	ID        int    `json:"id"`
	USER_ID   int    `json:"user_id"`
	Caption   string `json:"caption"`
	Image_url string `json:"image_url"` // maybe not this or update in order to have more photos, we ll see
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
