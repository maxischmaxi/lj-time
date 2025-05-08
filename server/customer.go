package main

import (
	"context"
	"time"

	"connectrpc.com/connect"
	customerv1 "github.com/maxischmaxi/ljtime-api/customer/v1"
	"github.com/maxischmaxi/ljtime-api/customer/v1/customerv1connect"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CustomerServer struct {
	customerv1connect.UnimplementedCustomerServiceHandler
	MongoClient *mongo.Client
	DBName string
}

type DbCustomer struct {
	Id bson.ObjectID `bson:"_id"`
	Name string `bson:"name"`
	Phone string `bson:"phone"`
	Email string `bson:"email"`
	Tag string `bson:"tag"`
	CreatedAt int64 `bson:"created_at"`
	UpdatedAt int64 `bson:"updated_at"`
}

func GetCustomerById(ctx context.Context, collection *mongo.Collection, id string) (*customerv1.Customer, error) {
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	var dbCustomer DbCustomer
	if err := collection.FindOne(ctx, bson.M{"_id": objId}).Decode(&dbCustomer); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}

		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &customerv1.Customer{
		Id: dbCustomer.Id.Hex(),
		Name: dbCustomer.Name,
		Phone: dbCustomer.Phone,
		Email: dbCustomer.Email,
		Tag: dbCustomer.Tag,
		CreatedAt: dbCustomer.CreatedAt,
		UpdatedAt: dbCustomer.UpdatedAt,
	}, nil
}

func (s *CustomerServer) GetCustomer(ctx context.Context, req *connect.Request[customerv1.GetCustomerRequest]) (*connect.Response[customerv1.GetCustomerResponse], error) {
	collection := s.MongoClient.Database(s.DBName).Collection("customers")
	id := req.Msg.Id

	customer, err := GetCustomerById(ctx, collection, id)
	if err != nil {
		return nil, err
	}

	res := connect.NewResponse(&customerv1.GetCustomerResponse{
		Customer: &customerv1.Customer{
			Id:       customer.Id,
			Name:     customer.Name,
			Phone:    customer.Phone,
			Email:    customer.Email,
			Tag:      customer.Tag,
			CreatedAt: customer.CreatedAt,
			UpdatedAt: customer.UpdatedAt,
		},
	})

	return res, nil
}

func (s *CustomerServer) CreateCustomer(ctx context.Context, req *connect.Request[customerv1.CreateCustomerRequest]) (*connect.Response[customerv1.CreateCustomerResponse], error) {
	collection := s.MongoClient.Database(s.DBName).Collection("customers")

	id := bson.NewObjectID()

	if _, err := collection.InsertOne(ctx, bson.M{
		"name": req.Msg.Customer.Name,
		"phone": req.Msg.Customer.Phone,
		"email": req.Msg.Customer.Email,
		"tag": req.Msg.Customer.Tag,
		"created_at": bson.Timestamp{
			T: uint32(time.Now().Unix()),
			I: 1,
		},
		"updated_at": bson.Timestamp{
			T: uint32(time.Now().Unix()),
			I: 1,
		},
		"_id": id,
	}); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&customerv1.CreateCustomerResponse{
		Customer: &customerv1.Customer{
			Id:    id.Hex(),
			Name:  req.Msg.Customer.Name,
			Phone: req.Msg.Customer.Phone,
			Email: req.Msg.Customer.Email,
			Tag:   req.Msg.Customer.Tag,
		},
	})

	return res, nil
}

func (s *CustomerServer) UpdateCustomer(ctx context.Context, req *connect.Request[customerv1.UpdateCustomerRequest]) (*connect.Response[customerv1.UpdateCustomerResponse], error) {
	collection := s.MongoClient.Database(s.DBName).Collection("customers")
	id := req.Msg.Customer.Id

	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	now := time.Now().Unix()

	if _, err := collection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{
		"$set": bson.M{
			"name": req.Msg.Customer.Name,
			"phone": req.Msg.Customer.Phone,
			"email": req.Msg.Customer.Email,
			"tag": req.Msg.Customer.Tag,
			"updated_at": now,
		},
	}); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&customerv1.UpdateCustomerResponse{
		Customer: &customerv1.Customer{
			Id:    id,
			Name:  req.Msg.Customer.Name,
			Phone: req.Msg.Customer.Phone,
			Email: req.Msg.Customer.Email,
			Tag:   req.Msg.Customer.Tag,
			CreatedAt: now,
			UpdatedAt: now,
		},
	})

	return res, nil
}

func (s *CustomerServer) DeleteCustomer(ctx context.Context, req *connect.Request[customerv1.DeleteCustomerRequest]) (*connect.Response[customerv1.DeleteCustomerResponse], error) {
	collection := s.MongoClient.Database(s.DBName).Collection("customers")
	id := req.Msg.Id

	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	if _, err := collection.DeleteOne(ctx, bson.M{"_id": objId }); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&customerv1.DeleteCustomerResponse{
		Id: id,
	})

	return res, nil
}

func (s *CustomerServer) GetCustomers(ctx context.Context, req *connect.Request[customerv1.GetCustomersRequest]) (*connect.Response[customerv1.GetCustomersResponse], error) {
	collection := s.MongoClient.Database(s.DBName).Collection("customers")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	defer cursor.Close(ctx)

	var customers []*customerv1.Customer

	for cursor.Next(ctx) {
		var dbCustomer DbCustomer
		if err := cursor.Decode(&dbCustomer); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		customers = append(customers, &customerv1.Customer{
			Id:       dbCustomer.Id.Hex(),
			Name:     dbCustomer.Name,
			Phone:    dbCustomer.Phone,
			Email:    dbCustomer.Email,
			Tag:      dbCustomer.Tag,
			CreatedAt: dbCustomer.CreatedAt,
			UpdatedAt: dbCustomer.UpdatedAt,
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&customerv1.GetCustomersResponse{
		Customers: customers,
	})

	return res, nil
}
