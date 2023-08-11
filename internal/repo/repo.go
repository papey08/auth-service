package repo

import (
	"auth-service/internal/app"
	"auth-service/internal/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// duplicateKeyErrorCode is a mongodb code of attempt of insertion of existing _id
const duplicateKeyErrorCode = 11000

type Repo struct {
	*mongo.Collection
}

type repoField struct {
	GUID         string `bson:"guid"`
	RefreshToken string `bson:"_id"`
}

func (r *Repo) InsertToken(ctx context.Context, user model.User, token string) error {
	_, err := r.InsertOne(ctx, repoField{
		GUID:         user.GUID,
		RefreshToken: token,
	})
	if we, ok := err.(mongo.WriteException); ok {
		for _, e := range we.WriteErrors {
			if e.Code == duplicateKeyErrorCode {
				return model.TokenCollisionError
			}
		}
	} else if err != nil {
		return model.RepoError
	}
	return nil
}

func (r *Repo) UpdateToken(ctx context.Context, oldToken, newToken string) error {
	res, err := r.UpdateOne(ctx,
		bson.M{"_id": oldToken},
		bson.M{"set": bson.M{"_id": newToken}})
	if res.ModifiedCount == 0 {
		return model.NoTokenError
	} else if err != nil {
		return model.RepoError
	}
	return nil
}

func (r *Repo) GetByRefreshToken(ctx context.Context, refreshToken string) (model.User, error) {
	var res repoField
	err := r.FindOne(ctx, bson.M{"_id": refreshToken}).Decode(&res)
	if err == mongo.ErrNoDocuments {
		return model.User{}, model.NoTokenError
	} else if err != nil {
		return model.User{}, model.RepoError
	}
	return model.User{
		GUID: res.GUID,
	}, nil
}

func New(c *mongo.Collection) app.Repo {
	return &Repo{
		Collection: c,
	}
}
