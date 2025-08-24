package dtos

// UserCreateDTO represents the data transfer object for creating a user
type UserCreateDTO struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

// UserUpdateDTO represents the data transfer object for updating a user
type UserUpdateDTO struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

// UserResponseDTO represents the data transfer object for user responses
type UserResponseDTO struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
