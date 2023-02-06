package projectEntity

type CreateProjectReq struct {
	ProjectId      	string     	`json:"project_id"`
	Title       	string     	`json:"title" validate:"required,min=3,max=20"`
	Color		 	string     	`json:"color" validate:"required,min=3"`
	UserId      	string     	`json:"user_id" validate:"required"`
	CreatedAt   	string     	`json:"created_at"`
	UpdatedAt   	string     	`json:"updated_at"`
}

type CreateProjectRes struct {
	ProjectId      	string     	`json:"project_id"`
	UserId      	string     	`json:"user_id" validate:"required"`
	Title       	string     	`json:"title" validate:"required,min=3"`
	Color		 	string     	`json:"color" validate:"required,min=3"`
}
