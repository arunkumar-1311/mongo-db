package repository

import (
	"context"

	"github.com/arunkumar-1311/mongo-db/dbconnection"
	"github.com/arunkumar-1311/mongo-db/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	usersCollection       = "user_details"
	departmentsCollection = "departments"
	rolesCollection       = "roles"
	ordersCollection      = "orders"
)

type Repo struct {
	DB *dbconnection.Database
}

type Repository interface {
	CreateUser(ctx context.Context, user models.User) (primitive.ObjectID, error)
	CreateDepartment(ctx context.Context, department models.Department) (primitive.ObjectID, error)
	CreateRole(ctx context.Context, role models.Role) (primitive.ObjectID, error)
	CreateOrders(ctx context.Context, orders []models.Order) ([]primitive.ObjectID, error)
	DeleteDepartment(ctx context.Context, id primitive.ObjectID) error
	DeleteRole(ctx context.Context, id primitive.ObjectID) error
	DeleteOrders(ctx context.Context, ids []primitive.ObjectID) error
}

func NewUserRepository() Repository {
	return Repo{DB: &dbconnection.Instance}
}

func (r Repo) CreateDepartment(ctx context.Context, department models.Department) (primitive.ObjectID, error) {
	department.Id = primitive.NewObjectID()
	_, err := r.DB.Db.Collection(departmentsCollection).InsertOne(ctx, department)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return department.Id, nil
}

func (r Repo) CreateRole(ctx context.Context, role models.Role) (primitive.ObjectID, error) {
	role.Id = primitive.NewObjectID()
	_, err := r.DB.Db.Collection(rolesCollection).InsertOne(ctx, role)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return role.Id, nil
}

func (r Repo) CreateOrders(ctx context.Context, orders []models.Order) ([]primitive.ObjectID, error) {
	if len(orders) == 0 {
		return nil, nil
	}

	docs := make([]interface{}, 0, len(orders))
	ids := make([]primitive.ObjectID, 0, len(orders))
	for i := range orders {
		orders[i].Id = primitive.NewObjectID()
		docs = append(docs, orders[i])
		ids = append(ids, orders[i].Id)
	}

	_, err := r.DB.Db.Collection(ordersCollection).InsertMany(ctx, docs)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func (r Repo) CreateUser(ctx context.Context, user models.User) (primitive.ObjectID, error) {
	user.Id = primitive.NewObjectID()
	_, err := r.DB.Db.Collection(usersCollection).InsertOne(ctx, user)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return user.Id, nil
}

func (r Repo) DeleteDepartment(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.DB.Db.Collection(departmentsCollection).DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r Repo) DeleteRole(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.DB.Db.Collection(rolesCollection).DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r Repo) DeleteOrders(ctx context.Context, ids []primitive.ObjectID) error {
	if len(ids) == 0 {
		return nil
	}
	_, err := r.DB.Db.Collection(ordersCollection).DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})
	return err
}
