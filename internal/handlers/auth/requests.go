package auth

type RequestLogin struct {
	Username string `json:"username"`
	Pass     string `json:"password"`
}

type RequestRegister struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Password2 string `json:"password2"`
	Email     string `json:"email"`
}

type RequestUserData struct {
	Username  string `json:"username"`
	Id        int64  `json:"id"`
	FirstName string `json:"first-name"`
	LastName  string `json:"last-name"`
}

type Validate interface {
	Validate() string
}
