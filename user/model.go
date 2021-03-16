package user

type LoginToken struct {
	RefreshToken   string
	Token          string
	TokenExpiresIn int64
	FirebaseToken  string
	TokenID        string
}

type User struct {
	ID int `gorm:"primary_key"`
}
