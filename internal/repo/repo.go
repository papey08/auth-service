package repo

import (
	"auth-service/internal/app"
	"auth-service/internal/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

const (
	// duplicateKeyErrorCode is a mongodb code of attempt of insertion of existing _id
	duplicateKeyErrorCode = 11000

	// removeExpTokensInterval is an interval between deletions of all expired refresh tokens from the database
	removeExpTokensInterval = time.Hour * 24
)

type Repo struct {
	*mongo.Collection
}

type repoField struct {
	RefreshToken string    `bson:"token"`
	GUID         string    `bson:"guid"`
	ExpiresAt    time.Time `bson:"expires_at"`
}

func (r *Repo) InsertToken(ctx context.Context, user model.User, token string, expiresAt time.Time) error {
	_, err := r.InsertOne(ctx, repoField{
		RefreshToken: token,
		GUID:         user.GUID,
		ExpiresAt:    expiresAt,
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

func (r *Repo) UpdateToken(ctx context.Context, oldToken, newToken string, expiresAt time.Time) error {
	filter := bson.M{"token": oldToken}
	update := bson.M{"$set": bson.M{
		"token":      newToken,
		"expires_at": expiresAt,
	}}
	res, err := r.UpdateOne(ctx, filter, update)
	if err != nil {
		return model.RepoError
	} else if res.ModifiedCount == 0 {
		return model.NoTokenError
	}
	return nil
}

func (r *Repo) GetByRefreshToken(ctx context.Context, refreshToken string) (model.User, time.Time, error) {
	var res repoField
	err := r.FindOne(ctx, bson.M{"token": refreshToken}).Decode(&res)
	if err == mongo.ErrNoDocuments {
		return model.User{}, time.Time{}, model.NoTokenError
	} else if err != nil {
		return model.User{}, time.Time{}, model.RepoError
	}
	return model.User{
		GUID: res.GUID,
	}, res.ExpiresAt, nil
}

func (r *Repo) RemoveExpiredTokens(ctx context.Context) {
	currTime := time.Now()
	_, err := r.DeleteMany(ctx, bson.M{"expires_at": bson.M{"$lt": currTime}})
	if err != nil {
		log.Println("removing of expired tokens error: ", err.Error())
	}
}

func New(ctx context.Context, c *mongo.Collection) app.Repo {
	r := Repo{
		Collection: c,
	}
	go func() {
		for {
			time.Sleep(removeExpTokensInterval)
			r.RemoveExpiredTokens(ctx)
		}
	}()
	return &r
}
