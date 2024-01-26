package models

import "time"

type Option struct {
	ID          string    `firestore:"id" json:"id"`
	Name        string    `firestore:"name" json:"name"`
	Choices     []Choice  `firestore:"choices" json:"choices"`
	Required    bool      `firestore:"required" json:"required"`
	MultiSelect bool      `firestore:"multi_select" json:"multi_select"`
	MaxSelect   int       `firestore:"max_select" json:"max_select"`
	CreatedAt   time.Time `firestore:"created_at" json:"created_at"`
	UpdatedAt   time.Time `firestore:"updated_at" json:"updated_at"`
	MenusId     []string  `firestore:"menus_id" json:"menus_id"`
	Menus       []Menu    `firestore:"menus" json:"menus"`
}

type Choice struct {
	ID     string  `firestore:"id" json:"id"`
	Index  int     `firestore:"index" json:"index"`
	Name   string  `firestore:"name" json:"name"`
	Price  float32 `firestore:"price" json:"price"`
	Prefix string  `firestore:"prefix" json:"prefix"`
}
