package models

import "time"

type Tables struct {
	ID         string    `firestore:"id" json:"id"`
	Name       string    `firestore:"name" json:"name"`
	Status     string    `firestore:"status" json:"status"`
	AccessUuid string    `firestore:"access_uuid" json:"access_uuid"`
	Orders     []Order   `firestore:"orders" json:"orders"`
	CreatedAt  time.Time `firestore:"created_at" json:"created_at"`
	UpdatedAt  time.Time `firestore:"updated_at" json:"updated_at"`
}

type Order struct {
	Uuid      string        `firestore:"uuid" json:"uuid"`
	MenuId    string        `firestore:"menu_id" json:"menu_id"`
	Name      string        `firestore:"name" json:"name"`
	Price     float32       `firestore:"price" json:"price"`
	Quantity  int           `firestore:"quantity" json:"quantity"`
	Status    string        `firestore:"status" json:"status"`
	Options   []OptionOrder `firestore:"options" json:"options"`
	CreatedAt time.Time     `firestore:"created_at" json:"created_at"`
}

type OptionOrder struct {
	Name  string   `firestore:"name" json:"name"`
	Value []string `firestore:"value" json:"value"`
	Price float32  `firestore:"price" json:"price"`
}

type OrderQuery struct {
	Orders []Order `json:"orders"`
}
