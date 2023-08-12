package repo

import (
	"auth-service/internal/app"
	"auth-service/internal/model"
	cryptotools "auth-service/pkg/crypto-tools"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

// removeExpTokensInterval is an interval between deletions of all expired refresh tokens from the database
const removeExpTokensInterval = time.Hour * 24

type Repo struct {
	*mongo.Collection
}

// repoField is a struct of record in database
type repoField struct {
	RefreshToken string    `bson:"token"`
	GUID         string    `bson:"guid"`
	ExpiresAt    time.Time `bson:"expires_at"`
}

func (r *Repo) InsertToken(ctx context.Context, user model.User, token string, expiresAt time.Time) error {
	// getting hash of refresh token
	encryptedToken, err := cryptotools.GenerateBcryptHash(token)
	if err != nil {
		return model.TokenCryptError
	}

	// inserting record with hashed token into database
	_, err = r.InsertOne(ctx, repoField{
		RefreshToken: encryptedToken,
		GUID:         user.GUID,
		ExpiresAt:    expiresAt,
	})

	if err != nil {
		return model.RepoError
	}
	return nil
}

func (r *Repo) UpdateToken(ctx context.Context, oldToken, newToken string, expiresAt time.Time) error {
	// finding hashed token by non-hashed token
	rf, err := r.getFieldByToken(ctx, oldToken)
	if err != nil {
		return err
	}

	// hashing new token
	encryptedNewToken, err := cryptotools.GenerateBcryptHash(newToken)
	if err != nil {
		return model.TokenCryptError
	}

	// updating hashed token in the database
	filter := bson.M{"token": rf.RefreshToken}
	update := bson.M{"$set": bson.M{
		"token":      encryptedNewToken,
		"expires_at": expiresAt,
	}}
	_, err = r.UpdateOne(ctx, filter, update)
	if err != nil {
		return model.RepoError
	}

	return nil
}

func (r *Repo) GetByRefreshToken(ctx context.Context, refreshToken string) (model.User, time.Time, error) {
	if rf, err := r.getFieldByToken(ctx, refreshToken); err != nil {
		return model.User{}, time.Time{}, err
	} else {
		return model.User{
			GUID: rf.GUID,
		}, rf.ExpiresAt, nil
	}
}

func (r *Repo) RemoveExpiredTokens(ctx context.Context) {
	currTime := time.Now()
	_, err := r.DeleteMany(ctx, bson.M{"expires_at": bson.M{"$lt": currTime}})
	if err != nil {
		log.Println("removing of expired tokens error: ", err.Error())
	}
}

// getFieldByToken search field in database where hashed token compares with non-hashed refreshToken
func (r *Repo) getFieldByToken(ctx context.Context, refreshToken string) (repoField, error) {
	cursor, err := r.Find(ctx, bson.D{})
	if err != nil {
		return repoField{}, model.RepoError
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

	for cursor.Next(ctx) {
		var f repoField
		if err = cursor.Decode(&f); err != nil {
			return repoField{}, model.RepoError
		}

		if cryptotools.CheckHash(f.RefreshToken, refreshToken) {
			return f, nil
		}
	}

	return repoField{}, model.NoTokenError

}

func New(ctx context.Context, c *mongo.Collection) app.Repo {
	r := Repo{
		Collection: c,
	}

	// once a day clears database from expired tokens
	go func() {
		for {
			time.Sleep(removeExpTokensInterval)
			r.RemoveExpiredTokens(ctx)
		}
	}()
	return &r
}
