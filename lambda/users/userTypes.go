package users

type User struct {
	UserID   uint   `json:"userID"`
	Username string `json:"userName"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
	IsAdmin  bool   `json:"isAdmin"`
}

