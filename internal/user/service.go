package user

import (
	"context"	
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}


func (s *Service) GetUser(ctx context.Context, id int64) (*User, error) {

	user, err := s.repo.GetByID(ctx,id)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) CreateUser(ctx context.Context, user *User) error {

	if err := validateUser(user); err != nil {
		return err
	}

	return s.repo.Create(ctx, user)
}

func (s *Service) UpdateUser(ctx context.Context, user *User) error {
	if err := validateUser(user); err != nil {
		return err
	}

	return s.repo.Update(ctx, user)
}


func validateUser(user *User) error {

	if user == nil {
		return ErrInvalidUser
	}

	if user.Name == "" {
		return ErrInvalidUser
	}
	if user.Email == "" {
		return ErrInvalidUser
	}
	return nil
}