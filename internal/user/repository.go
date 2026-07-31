package user

import (
	"context"
	"sync"
)

type Repository interface {
	GetByID(ctx context.Context, id int64) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
}

type MemoryRepository struct {
	mu    sync.RWMutex
	users map[int64]*User
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		users: make(map[int64]*User),
	}
}

func (r *MemoryRepository) Create(ctx context.Context, user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; exists {
		return ErrUserAlreadyExists
	}

	r.users[user.ID] = user
	return nil
}

func (r *MemoryRepository) GetByID(ctx context.Context, id int64) (*User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if user, exists := r.users[id]; exists {
		return user, nil
	}

	return nil, ErrUserNotFound
}

func (r *MemoryRepository) Update(ctx context.Context, user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return ErrUserNotFound
	}

	r.users[user.ID] = user
	return nil
}

