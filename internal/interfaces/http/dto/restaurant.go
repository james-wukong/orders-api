package dto

// CreateRestaurantRequest is what the client sends (POST /register)
type CreateRestaurantRequest struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	Description string `json:"description"`
	Phone       string `json:"phone"`
	Email       string `json:"email" binding:"required,email"`
	Address     string `json:"address"`
	City        string `json:"city"`
	State       string `json:"state"`
}

// RestaurantResponse is what we send back to the client
type RestaurantResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Email       string `json:"email,omitempty"`
	Address     string `json:"address,omitempty"`
	City        string `json:"city,omitempty"`
	State       string `json:"state,omitempty"`
	CreatedAt   string `json:"created_at"`
}
