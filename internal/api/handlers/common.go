package handlers

// BulkDeleteRequest defines the payload for bulk-deleting items.
type BulkDeleteRequest struct {
	IDs []int64 `json:"ids" binding:"required,min=1"`
}

// PaginationMeta contains metadata for paginated responses.
type PaginationMeta struct {
	TotalItems   int `json:"total_items"`
	TotalPages   int `json:"total_pages"`
	CurrentPage  int `json:"current_page"`
	ItemsPerPage int `json:"items_per_page"`
}