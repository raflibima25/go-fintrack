package response

type Pagination struct {
	CurrentPage int   `json:"current_page"`
	TotalPage   int   `json:"total_page"`
	TotalItems  int64 `json:"total_items"`
	ItemPerPage int   `json:"item_per_page"`
}
