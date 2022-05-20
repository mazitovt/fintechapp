package server

type ErrMissingArgument struct {
	name string
}

type ErrWrongHttpMethod struct {
	used   string
	should string
}

type ErrNoSuchPath struct {
	path string
}

type ErrInternalError struct {
}

// ------SignUp------

type ErrEmailAlreadyUsed struct {
	// TODO maybe not disclose user email?
	email string
}

type ErrPasswordTooShort struct {
	// should be at least 8 symbols
}

type ErrEmailTooShort struct {
	// should be at least 8 symbols
}

type ErrInvalidEmail struct {
	// didn't pass regex check
}

// ------SignIn------

type ErrInvalidCredentials struct {
	// email and password doesn't match
}

// ------Refresh/Parse------

type ErrInvalidToken struct {
	// received token is invalid
	// access: expired or invalid
	// refresh: invalid, expired or not found in db
}

//func NewErrMissingArgument(args ...string) error {
//
//}
//
//func (e *ErrMissingArgument) Error() string {
//	return fmt.Sprintf()
//}
