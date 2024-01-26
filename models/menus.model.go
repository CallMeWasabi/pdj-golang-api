package models

import (
	"time"
)

type Menu struct {
	ID         string    `firestore:"id" json:"id"`
	Name       string    `firestore:"name" json:"name"`
	Price      float32   `firestore:"price" json:"price"`
	Status     string    `firestore:"status" json:"status"`
	CreatedAt  time.Time `firestore:"created_at" json:"created_at"`
	UpdatedAt  time.Time `firestore:"updated_at" json:"updated_at"`
	MenuTypeId string    `firestore:"menu_type_id" json:"menu_type_id"`
	OptionsId  []string  `firestore:"options_id" json:"options_id"`
	Options    []Option  `firestore:"options" json:"options"`
}
