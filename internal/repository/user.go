package repository

import (
	"backendTemp/internal/model"
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
}

type MongoUserRepo struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &MongoUserRepo{db: db}
}

func (r *MongoUserRepo) Create(ctx context.Context, user *model.User) error {
	_, err := r.db.Collection("users").InsertOne(ctx, user)
	return err
}

func (r *MongoUserRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.Collection("users").FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
