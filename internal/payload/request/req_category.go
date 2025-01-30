package request

type CategoryRequest struct {
	Name      string `json:"name" binding:"required"`
	Color     string `json:"color"`
	IconColor string `json:"icon_color"`
}

type UpdateCategoryRequest struct {
	Name      string `json:"name" binding:"required"`
	Color     string `json:"color"`
	IconColor string `json:"icon_color"`
}
