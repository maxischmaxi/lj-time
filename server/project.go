package main

import (
	"context"
	"time"

	"connectrpc.com/connect"
	projectv1 "github.com/maxischmaxi/ljtime-api/project/v1"
	"github.com/maxischmaxi/ljtime-api/project/v1/projectv1connect"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ProjectServer struct {
	projectv1connect.UnimplementedProjectServiceHandler
	MongoClient *mongo.Client
	DBName string
}

type DbProject struct {
	Id bson.ObjectID `bson:"_id"`
	Name string `bson:"name"`
	Description string `bson:"description"`
	CustomerId bson.ObjectID `bson:"customer_id"`
	CreatedAt int64 `bson:"created_at"`
	UpdatedAt int64 `bson:"updated_at"`
}

func GetProjectById(ctx context.Context, collection *mongo.Collection, id string) (*projectv1.Project, error) {
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	var dbProject DbProject
	if err := collection.FindOne(ctx, bson.M{"_id": objId}).Decode(&dbProject); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}

		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return &projectv1.Project{
		Id: dbProject.Id.Hex(),
		Name: dbProject.Name,
		Description: dbProject.Description,
		CustomerId: dbProject.CustomerId.Hex(),
		CreatedAt: dbProject.CreatedAt,
		UpdatedAt: dbProject.UpdatedAt,
	}, nil
}

func (s *ProjectServer) GetProject(ctx context.Context, req *connect.Request[projectv1.GetProjectRequest]) (*connect.Response[projectv1.GetProjectResponse], error) {
	collection := s.MongoClient.Database(s.DBName).Collection("projects")
	id := req.Msg.Id

	project, err := GetProjectById(ctx, collection, id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&projectv1.GetProjectResponse{
		Project: project,
	})

	return res, nil
}

func (s *ProjectServer) CreateProject(ctx context.Context, req *connect.Request[projectv1.CreateProjectRequest]) (*connect.Response[projectv1.CreateProjectResponse], error) {
	collection := s.MongoClient.Database(s.DBName).Collection("projects")

	id := bson.NewObjectID()
	customerId, err := bson.ObjectIDFromHex(req.Msg.Project.CustomerId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	customer, err := GetCustomerById(ctx, collection, customerId.Hex())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	now := time.Now().Unix()

	if _, err := collection.InsertOne(ctx, bson.M{
		"name": req.Msg.Project.Name,
		"description": req.Msg.Project.Description,
		"customer_id": customer.Id,
		"created_at": now,
		"updated_at": now,
	}); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&projectv1.CreateProjectResponse{
		Project: &projectv1.Project{
			Id:    id.Hex(),
			Name:  req.Msg.Project.Name,
			Description: req.Msg.Project.Description,
			CustomerId: customer.Id,
			CreatedAt: now,
			UpdatedAt: now,
		},
	})

	return res, nil
}

func (s *ProjectServer) UpdateProject(ctx context.Context, req *connect.Request[projectv1.UpdateProjectRequest]) (*connect.Response[projectv1.UpdateProjectResponse], error) {
	collection := s.MongoClient.Database(s.DBName).Collection("projects")

	objId, err := bson.ObjectIDFromHex(req.Msg.Project.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	customerId, err := bson.ObjectIDFromHex(req.Msg.Project.CustomerId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	customer, err := GetCustomerById(ctx, collection, customerId.Hex())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	now := time.Now().Unix()

	if _, err := collection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{
		"$set": bson.M{
			"name": req.Msg.Project.Name,
			"description": req.Msg.Project.Description,
			"customer_id": customer.Id,
			"updated_at": now,
		},
	}); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&projectv1.UpdateProjectResponse{
		Project: &projectv1.Project{
			Id:    req.Msg.Project.Id,
			Name:  req.Msg.Project.Name,
			Description: req.Msg.Project.Description,
			CustomerId: customer.Id,
			CreatedAt: now,
			UpdatedAt: now,
		},
	})

	return res, nil
}

func (s *ProjectServer) DeleteProject(ctx context.Context, req *connect.Request[projectv1.DeleteProjectRequest]) (*connect.Response[projectv1.DeleteProjectResponse], error) {
	collection := s.MongoClient.Database(s.DBName).Collection("projects")

	objId, err := bson.ObjectIDFromHex(req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	if _, err := collection.DeleteOne(ctx, bson.M{"_id": objId}); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&projectv1.DeleteProjectResponse{
		Id: req.Msg.Id,
	})

	return res, nil
}

func (s *ProjectServer) GetProjects(ctx context.Context, req *connect.Request[projectv1.GetProjectsRequest]) (*connect.Response[projectv1.GetProjectsResponse], error) {
	collection := s.MongoClient.Database(s.DBName).Collection("projects")
	var dbProjects []DbProject
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var dbProject DbProject
		if err := cursor.Decode(&dbProject); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		dbProjects = append(dbProjects, dbProject)
	}

	if err := cursor.Err(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	var projects []*projectv1.Project
	for _, dbProject := range dbProjects {
		projects = append(projects, &projectv1.Project{
			Id: dbProject.Id.Hex(),
			Name: dbProject.Name,
			Description: dbProject.Description,
			CustomerId: dbProject.CustomerId.Hex(),
			CreatedAt: dbProject.CreatedAt,
			UpdatedAt: dbProject.UpdatedAt,
		})
	}

	res := connect.NewResponse(&projectv1.GetProjectsResponse{
		Projects: projects,
	})

	return res, nil
}
