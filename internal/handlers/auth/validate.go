package auth

import "fmt"

func (u *RequestUserData) Validate() string {
	switch {
	case len(u.Username) == 0:
		return fmt.Sprintf("empty username")
	case len(u.FirstName) == 0:
		return fmt.Sprintf("empty first name")
	case len(u.LastName) == 0:
		return fmt.Sprintf("empty last name")
	}
	return ""
}

func (u *RequestLogin) Validate() string {
	switch {
	case len(u.Username) == 0:
		return fmt.Sprintf("empty username")
	case len(u.Pass) == 0:
		return fmt.Sprintf("empty password")
	}
	return ""
}

func (r *RequestRegister) Validate() string {
	switch {
	case len(r.Username) == 0:
		return fmt.Sprintf("empty username")
	case len(r.Password) == 0:
		if len(r.Password) < 6 {
			return fmt.Sprintf("The password length is less than 6")
		}
		return fmt.Sprintf("empty password")
	case len(r.Password2) == 0:
		return fmt.Sprintf("empty password2")
	case r.Password != r.Password2:
		return fmt.Sprintf("invalid password")
	case len(r.Email) == 0:
		return fmt.Sprintf("empty email")
	}
	return ""
}

func ValidateEmail(email string) bool {
	panic(email)
}
