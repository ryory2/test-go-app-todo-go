package response

type SuccessResponse struct {
	Status       int         `json:"status"`
	TotalResults int         `json:"totalResults,omitempty"`
	Resources    interface{} `json:"Resources,omitempty"`
	ID           uint        `json:"id,omitempty"`
	Title        string      `json:"title,omitempty"`
	Description  string      `json:"description,omitempty"`
	DueDate      string      `json:"due_date,omitempty"`
	IsCompleted  bool        `json:"is_completed,omitempty"`
	CreatedAt    string      `json:"created_at,omitempty"`
	UpdatedAt    string      `json:"updated_at,omitempty"`
}

type ErrorResponse struct {
	Status string `json:"status"`
	Detail string `json:"detail"`
}
