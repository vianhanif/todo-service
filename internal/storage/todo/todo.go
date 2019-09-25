// template version: 1.0.9
package todo

import (
	"time"
)

// Todo ,
type Todo struct {
	// ID ,
	ID int `json:"id"`
	// Title ,
	Title string `json:"title"`
	// Detail ,
	Detail string `json:"detail"`
	// CreatedAt ,
	CreatedAt time.Time `json:"createdAt"`
	// UpdatedAt ,
	UpdatedAt time.Time `json:"updatedAt"`
	// DeletedAt ,
	DeletedAt *time.Time `json:"deletedAt"`
}
