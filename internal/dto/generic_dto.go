package dto

// ====================================================================
// ============================ Common DTO ============================
// ====================================================================

// Filter struct to describe filtering and pagination parameters for generic requests.
// @Description Generic filter structure for pagination, searching, and ordering
// @Page Page is the current page number for pagination (optional, starts from 1)
// @PerPage PerPage is the number of items to display per page (optional)
// @Keyword Keyword is used for searching/filtering records by text content
// @OrderBy OrderBy specifies the field to sort by, prefix with '-' for descending order
// @Tags Request Filters
type Filter struct {
	Page    int    `json:"page" example:"1" validate:"number" doc:"Current page number for pagination"`
	PerPage int    `json:"per_page" example:"10" validate:"number" doc:"Number of items to display per page"`
	Keyword string `json:"keyword" example:"search term" validate:"" doc:"Search keyword for filtering records"`
	OrderBy string `json:"order_by" example:"-created_at" validate:"" doc:"Field to order by, prefix with '-' for descending order"`
}
