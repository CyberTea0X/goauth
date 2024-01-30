package controllers

func ErrToMap(err error) map[string]interface{} {
	return map[string]interface{}{
		"error": err.Error(),
	}
}

type ErrNoTokenSpecified struct{}

func (e ErrNoTokenSpecified) Error() string {
	return "no token specified"
}

type ErrInvalidAccessToken struct{}

func (e ErrInvalidAccessToken) Error() string {
	return "invalid access token"
}

type ErrTokenExpired struct{}

func (e ErrTokenExpired) Error() string {
	return "token expired"
}
