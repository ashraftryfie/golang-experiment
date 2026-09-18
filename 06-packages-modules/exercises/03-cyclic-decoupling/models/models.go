package models

type User struct {
	ID   string
	Name string
}

type Order struct {
	ID     string
	UserID string
	Amount float64
}
