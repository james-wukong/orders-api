package restaurant

import "errors"

var (
	ErrRestaurantNotFound = errors.New("restaurant not found")
	ErrSlugAlreadyUsed    = errors.New("slug already used")
)
