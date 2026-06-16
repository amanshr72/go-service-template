package user

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoRepo struct {
	col *mongo.Collection
}

func NewMongoRepository(col *mongo.Collection) Repository {
	return &mongoRepo{col: col}
}

func MongoMigrate(col *mongo.Collection) error {
	_, err := col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	return err
}

func (r *mongoRepo) Create(u *User) error {
	count, err := r.col.CountDocuments(context.Background(), bson.D{})
	if err != nil {
		return err
	}
	u.ID = int(count) + 1
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	_, err = r.col.InsertOne(context.Background(), u)
	return err
}

func (r *mongoRepo) FindAll() ([]User, error) {
	cursor, err := r.col.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	var users []User
	if err := cursor.All(context.Background(), &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *mongoRepo) FindByID(id int) (*User, error) {
	var u User
	err := r.col.FindOne(context.Background(), bson.M{"id": id}).Decode(&u)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.New("user not found")
	}
	return &u, err
}

func (r *mongoRepo) FindByActive(active bool) ([]User, error) {
	cursor, err := r.col.Find(context.Background(), bson.M{"is_active": active})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	var users []User
	if err := cursor.All(context.Background(), &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *mongoRepo) Count() (int, error) {
	count, err := r.col.CountDocuments(context.Background(), bson.D{})
	return int(count), err
}

func (r *mongoRepo) Update(u *User) error {
	u.UpdatedAt = time.Now()
	result, err := r.col.UpdateOne(
		context.Background(),
		bson.M{"id": u.ID},
		bson.M{"$set": u},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *mongoRepo) Delete(id int) error {
	result, err := r.col.DeleteOne(context.Background(), bson.M{"id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("user not found")
	}
	return nil
}
