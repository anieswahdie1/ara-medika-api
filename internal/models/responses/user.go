package responses

import (
	"time"

	"github.com/anieswahdie1/ara-medika-api.git/internal/models/entities"
)

type UserResponse struct {
	ID         uint             `json:"id"`
	Name       string           `json:"name"`
	Email      string           `json:"email"`
	Role       string           `json:"role"`
	AccessMenu []entities.Menus `json:"access_menus"`
	CreatedAt  time.Time        `json:"created_at"`
}

type GetUsers struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Status string `json:"status"`
}
