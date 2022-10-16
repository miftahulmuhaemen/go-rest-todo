package user

import (
	core "go-rest-todo/core/user"
)

type Repository interface {
	Create(core.User) (core.User, error)
}
