package dto

type Image struct {
	Thumbnail string `json:"thumbnail"`
	Mobile    string `json:"mobile"`
	Tablet    string `json:"tablet"`
	Desktop   string `json:"desktop"`
}
type ProductResponse struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Cuisines    []string `json:"cuisines,omitempty"`
	Category    string   `json:"category"`
	Price       float64  `json:"price"`
	Description string   `json:"description,omitempty"`
	Rating      float32  `json:"rating"`
	Image       []Image  `json:"image"`
}
