package repository

import (
	"context"

	"github.com/arunkumar-1311/mongo-db/dbconnection"
	"github.com/arunkumar-1311/mongo-db/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	GetUser(ctx context.Context, id primitive.ObjectID) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
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

func (r Repo) GetUser(ctx context.Context, id primitive.ObjectID) (models.User, error) {

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"_id": id}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         departmentsCollection,
			"localField":   "department_id",
			"foreignField": "_id",
			"as":           "department",
		}}},
		{{Key: "$unwind", Value: bson.M{
			"path":                       "$department",
			"preserveNullAndEmptyArrays": true,
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         rolesCollection,
			"localField":   "role_id",
			"foreignField": "_id",
			"as":           "role",
		}}},
		{{Key: "$unwind", Value: bson.M{
			"path":                       "$role",
			"preserveNullAndEmptyArrays": true,
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         ordersCollection,
			"localField":   "order_ids",
			"foreignField": "_id",
			"as":           "orders",
		}}},
	}

	cursor, err := r.DB.Db.Collection(usersCollection).Aggregate(ctx, pipeline)
	if err != nil {
		return models.User{}, err
	}
	defer cursor.Close(ctx)

	var results []models.User
	if err := cursor.All(ctx, &results); err != nil {
		return models.User{}, err
	}

	if len(results) == 0 {
		return models.User{}, err
	}

	return results[0], nil
}

func (r Repo) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var result models.User
	err := r.DB.Db.Collection(usersCollection).FindOne(ctx, bson.M{"email": email}).Decode(&result)
	if err != nil {
		return models.User{}, err
	}
	return result, nil
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
