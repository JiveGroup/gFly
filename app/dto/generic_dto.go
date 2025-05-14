package dto

type Filter struct {
	Page    int    `json:"page" example:"1" validate:"number" doc:"Current page number"`
	PerPage int    `json:"per_page" example:"10" validate:"number" doc:"Number of items per page"`
	Keyword string `json:"keyword" example:"" validate:"" doc:"Search keyword"`
	OrderBy string `json:"order_by" example:"-full_name" validate:"" doc:"Field to order by, prefix with '-' for descending"`
}

type Meta struct {
	Page    int `json:"page,omitempty" example:"1" doc:"Current page number"`
	PerPage int `json:"per_page,omitempty" example:"10" doc:"Number of items per page"`
	Total   int `json:"total" example:"1354" doc:"Total number of records"`
}
