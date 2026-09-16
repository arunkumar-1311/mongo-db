package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id           primitive.ObjectID `json:"id" bson:"_id"`
	Name         string             `json:"name" binding:"required"`
	Salary       int64              `json:"salary" binding:"required"`
	Address      `json:"address"`
	Department   Department           `json:"department" binding:"required"`
	DepartmentId primitive.ObjectID   `json:"-" bson:"department_id"`
	Role         Role                 `json:"role" binding:"required"`
	RoleId       primitive.ObjectID   `json:"-" bson:"role_id"`
	Orders       []Order              `json:"orders"`
	OrderIds     []primitive.ObjectID `json:"-" bson:"order_ids"`
}

type Address struct {
	Id      int    `json:"id"`
	Street  string `json:"street"`
	City    string `json:"city"`
	Country string `json:"country"`
}

// Department is stored in its own "departments" collection and referenced
// from User via DepartmentId.
type Department struct {
	Id    primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Name  string             `json:"name" binding:"required"`
	Code  string             `json:"code" binding:"required"`
	Floor int                `json:"floor"`
}

// Role is stored in its own "roles" collection and referenced from User
// via RoleId.
type Role struct {
	Id          primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Title       string             `json:"title" binding:"required"`
	Level       string             `json:"level" binding:"required"`
	Permissions []string           `json:"permissions"`
}

// Order is stored in its own "orders" collection. A user can have many
// orders, referenced from User via OrderIds.
type Order struct {
	Id       primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Product  string             `json:"product" binding:"required"`
	Quantity int                `json:"quantity" binding:"required"`
	Price    float64            `json:"price" binding:"required"`
}
