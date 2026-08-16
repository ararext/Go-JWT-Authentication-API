package mock

import (
	"context"

	"github.com/ararext/Go-JWT-Authentication-API/internal/models"
	"github.com/ararext/Go-JWT-Authentication-API/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserRepository is an in-memory implementation of repository.UserRepository,
// used only in tests — never wired into the real application.
type UserRepository struct {
	usersByEmail map[string]*models.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		usersByEmail: make(map[string]*models.User),
	}
}

func (m *UserRepository) Create(ctx context.Context, user *models.User) error {
	if user.ID.IsZero() {
		user.ID = primitive.NewObjectID()
	}
	m.usersByEmail[user.Email] = user
	return nil
}

func (m *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	user, ok := m.usersByEmail[email]
	if !ok {
		return nil, repository.ErrUserNotFound
	}
	return user, nil
}

func (m *UserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	for _, u := range m.usersByEmail {
		if u.ID.Hex() == id {
			return u, nil
		}
	}
	return nil, repository.ErrUserNotFound
}
