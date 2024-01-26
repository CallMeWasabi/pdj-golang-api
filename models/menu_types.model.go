package models

import "time"

type MenuType struct {
	ID        string    `firestore:"id" json:"id"`
	Name      string    `firestore:"name" json:"name"`
	Status    string    `firestore:"status" json:"status"`
	MenusId   []string  `firestore:"menus_id" json:"menus_id"`
	Menus     []Menu    `firestore:"menus" json:"menus"`
	CreatedAt time.Time `firestore:"created_at" json:"created_at"`
	UpdatedAt time.Time `firestore:"updated_at" json:"updated_at"`
}
