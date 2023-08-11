package repo

import (
	"auth-service/internal/app"
	"auth-service/internal/model"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Repo struct {
	*mongo.Collection
}

func (r *Repo) SetSession(ctx context.Context, user model.User, refreshToken string, expiresAt time.Time) error {
	// TODO: implement
	return nil
}

func (r *Repo) GetByRefreshToken(ctx context.Context, refreshToken string) (model.User, error) {
	// TODO: implement
	return model.User{}, nil
}

func New(c *mongo.Collection) app.Repo {
	return &Repo{
		Collection: c,
	}
}
