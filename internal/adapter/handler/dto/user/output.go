package user

import domain "github.com/MathieuRocher/amazing_domain"

type UserOutput struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func FromDomain(u *domain.User) *UserOutput {
	return &UserOutput{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
		Role:  u.Role.String(), // ex: "Admin"
	}
}

func ListFromDomain(users []domain.User) []UserOutput {
	result := make([]UserOutput, len(users))
	for i, u := range users {
		result[i] = *FromDomain(&u)
	}
	return result
}
