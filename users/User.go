package users

type User struct {
	ID           int     `table:"users" json:"id"`
	Username     string  `json:"username"`
	Password     string  `json:"-"`
	Email        string  `json:"email"`
	MobileNumber *string `json:"mobileNumber"`
}
