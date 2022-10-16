package user

import (
	core "go-rest-todo/core/user"
)

type Service interface {
	Register(core.User) (core.User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) Register(user core.User) (core.User, error) {
	return core.User{}, nil
}
