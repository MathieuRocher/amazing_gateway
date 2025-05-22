package user

type UpdateUserInput struct {
	Name  *string `form:"name" validate:"omitempty,min=2"`
	Email *string `form:"email" validate:"omitempty,email"`
	Role  *string `form:"role" validate:"omitempty,oneof=Trainee Trainer Admin"`
}
