package AuthStructs

import "time"

type User struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	Deleted   bool
	Online    bool

	Username string `gorm:"unique"`
	Password string
	Email    string `gorm:"unique"`
	Role     int

	EmailConfirmed        bool
	EmailVerificationCode string

	CanChangePassword  bool
	PasswordChangeCode string
}
