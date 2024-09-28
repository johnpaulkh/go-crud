package repository

import (
	"context"
	"johnpaulkh/go-crud/api/model"
)

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
