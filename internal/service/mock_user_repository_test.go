package service

import (
	"context"

	"github.com/ararext/Go-JWT-Authentication-API/internal/models"
	"github.com/ararext/Go-JWT-Authentication-API/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockUserRepository struct {
	usersByEmail map[string]*models.User
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		usersByEmail: make(map[string]*models.User),
	}
}

func (m *mockUserRepository) Create(ctx context.Context, user *models.User) error {
	user.ID = primitive.NewObjectID()
	m.usersByEmail[user.Email] = user
	return nil
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	user, ok := m.usersByEmail[email]
	if !ok {
		return nil, repository.ErrUserNotFound
	}
	return user, nil
}

func (m *mockUserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	for _, u := range m.usersByEmail {
		if u.ID.Hex() == id {
			return u, nil
		}
	}
	return nil, repository.ErrUserNotFound
}