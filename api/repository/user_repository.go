package repository

import (
	"context"
	"johnpaulkh/go-crud/api/config"
	"johnpaulkh/go-crud/api/model"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepository struct {
	client *mongo.Client
	config *config.Configuration
}

type UserRepository interface {
	Create(
		request model.User,
		ctx context.Context) (*model.User, error)
	Update(
		id string,
		request model.User,
		ctx context.Context) (*model.User, error)
	Get(
		id string,
		ctx context.Context) (*model.User, error)
	List(
		page int,
		size int,
		ctx context.Context) (*model.Page[model.User], error)
}

func NewUserRepository(client *mongo.Client, config *config.Configuration) UserRepository {
	return &userRepository{
		client: client,
		config: config,
	}
}

func getCollection(app *userRepository) *mongo.Collection {
	return app.client.Database(app.config.Database.DbName).Collection(app.config.Database.Collection)
}

func (app *userRepository) Create(user model.User, ctx context.Context) (*model.User, error) {
	collection := getCollection(app)

	insertResult, err := collection.InsertOne(ctx, user)

	if oidResult, ok := insertResult.InsertedID.(primitive.ObjectID); ok {
		filter := bson.D{{Key: "_id", Value: oidResult}}
		var result model.User
		collection.FindOne(ctx, filter).Decode(&result)
		return &result, nil
	} else {
		return nil, err
	}
}

func (app *userRepository) Update(id string, user model.User, ctx context.Context) (*model.User, error) {
	collection := getCollection(app)

	update := bson.D{{
		Key: "$set", Value: bson.D{
			{Key: "name", Value: user.Name},
			{Key: "email", Value: user.Email},
			{Key: "phone", Value: user.Phone},
		},
	}}

	oid, _ := primitive.ObjectIDFromHex(id)
	updateResult, err := collection.UpdateByID(ctx, oid, update)

	if count := updateResult.MatchedCount; count == 1 {
		return &user, nil
	} else {
		return nil, err
	}
}

func (app *userRepository) Get(id string, ctx context.Context) (*model.User, error) {
	collection := getCollection(app)

	oId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{Key: "_id", Value: oId}}

	var result model.User
	collection.FindOne(ctx, filter).Decode(&result)

	return &result, nil
}

func (app *userRepository) List(page int, size int, ctx context.Context) (*model.Page[model.User], error) {

	skip := page * size
	opts := options.Find()
	opts.SetLimit(int64(size))
	opts.SetSkip(int64(skip))

	collection := getCollection(app)

	cursor, err := collection.Find(ctx, bson.D{{}}, opts)
	if err != nil {
		logrus.Error("Error during listing users in repo ", err)
		return nil, err
	}

	count, err := collection.CountDocuments(ctx, bson.D{{}})
	if err != nil {
		logrus.Error("Error during counting users in repo ", err)
		return nil, err
	}

	var users []*model.User
	for cursor.Next(ctx) {
		var user model.User
		if err := cursor.Decode(&user); err != nil {
			logrus.Error("Error during decoding user during listing process ", err)
			return nil, err
		}

		users = append(users, &user)
	}

	cursor.Close(ctx)

	result := model.Page[model.User]{
		Page:    page,
		Count:   count,
		Content: users,
	}
	return &result, nil
}
