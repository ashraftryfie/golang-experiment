package main

import "fmt"

type BaseEntity struct {
	ID        string
	CreatedAt int64
}

func (b BaseEntity) Summary() string {
	return fmt.Sprintf("ID: %s (Created: %d)", b.ID, b.CreatedAt)
}

// User embeds BaseEntity (composition, not inheritance)
type User struct {
	BaseEntity // Promoted fields: User.ID, User.CreatedAt, User.Summary()
	Username   string
	Email      string
}

// Admin embeds User (multi-level composition)
type Admin struct {
	User
	Role string
}

// Overriding promoted method on Admin
func (a Admin) Summary() string {
	return fmt.Sprintf("Admin: %s, Role: %s, %s", a.Username, a.Role, a.BaseEntity.Summary())
}

func main() {
	u := User{
		BaseEntity: BaseEntity{ID: "usr_123", CreatedAt: 1710000000},
		Username:   "ashraf",
		Email:      "ashraf@example.com",
	}

	// Promoted fields accessed directly
	fmt.Printf("User ID: %s\n", u.ID)
	fmt.Printf("User Summary: %s\n", u.Summary())

	a := Admin{
		User: User{
			BaseEntity: BaseEntity{ID: "adm_999", CreatedAt: 1710005000},
			Username:   "superadmin",
			Email:      "admin@system.io",
		},
		Role: "PlatformOps",
	}

	// Uses Admin's overridden Summary()
	fmt.Println(a.Summary())
}
