package dto

type ProductFilter struct {
	Name      string  `query:"name"`      // partial, case-insensitive
	Category  string  `query:"category"`  // exact match
	MinPrice  float64 `query:"minPrice"`  // inclusive lower bound (0 = no limit)
	MaxPrice  float64 `query:"maxPrice"`  // inclusive upper bound (0 = no limit)
	SortBy    string  `query:"sortBy"`    // "name" | "price" | "rating"  (default: "name")
	SortOrder string  `query:"sortOrder"` // "asc" | "desc"               (default: "asc")
}

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
