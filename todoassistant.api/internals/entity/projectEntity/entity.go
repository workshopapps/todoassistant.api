package projectEntity

type CreateProjectReq struct {
	ProjectId string `json:"project_id"`
	Title     string `json:"title" validate:"required,min=3,max=20"`
	Color     string `json:"color" validate:"required,min=3"`
	UserId    string `json:"user_id" validate:"required"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CreateProjectRes struct {
	ProjectId string `json:"project_id"`
	UserId    string `json:"user_id" validate:"required"`
	Title     string `json:"title" validate:"required,min=3"`
	Color     string `json:"color" validate:"required,min=3"`
}

type EditProjectReq struct {
	ProjectId string `json:"project_id"`
	Title     string `json:"title"`
	Color     string `json:"color"`
	UserId    string `json:"user_id"`
	UpdatedAt string `json:"updated_at"`
}

type EditProjectRes struct {
	ProjectId string `json:"project_id"`
	UserId    string `json:"user_id"`
	Title     string `json:"title"`
	Color     string `json:"color"`
	UpdatedAt string `json:"updated_at"`
}

type GetProjectRes struct {
	ProjectId string `json:"project_id"`
	Title     string `json:"title"`
	Color     string `json:"color"`
	UserId    string `json:"user_id"`
}
