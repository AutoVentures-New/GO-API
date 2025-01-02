package auth

import "errors"

var ErrUserAlreadyExists = errors.New("User already exists")

var ErrCompanyAlreadyExists = errors.New("Company already exists")
