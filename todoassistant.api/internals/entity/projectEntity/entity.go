package projectEntity

type CreateProjectReq struct {
	ProjectId      	string     	`json:"project_id"`
	UserId      	string     	`json:"user_id" validate:"required"`
	Title       	string     	`json:"title" validate:"required,min=3"`
	Color		 	string     	`json:"color" validate:"required,min=3"`
	CreatedAt   	string     	`json:"created_at"`
	UpdatedAt   	string     	`json:"updated_at"`
	IsDeleted		string		`json:"is_deleted"`
}

type CreateProjectRes struct {
	ProjectId      	string     	`json:"project_id"`
	UserId      	string     	`json:"user_id" validate:"required"`
	Title       	string     	`json:"title" validate:"required,min=3"`
	Color		 	string     	`json:"color" validate:"required,min=3"`
	CreatedAt   	string     	`json:"created_at"`
	UpdatedAt   	string     	`json:"updated_at"`
	IsDeleted		string		`json:"is_deleted"`
}
