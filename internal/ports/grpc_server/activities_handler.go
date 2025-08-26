package handler

import (
	"context"
	"crm/api/proto/pb"
	"crm/internal/adapters/database/db"
	"crm/internal/core/services"
	"database/sql"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ActivityHandler struct {
	activityService services.ActivityService
	pb.UnimplementedActivityServiceServer
}

func NewActivityHandler(service services.ActivityService) *ActivityHandler {
	return &ActivityHandler{activityService: service}
}

func (h *ActivityHandler) CreateActivity(ctx context.Context, req *pb.CreateActivityRequest) (*pb.CreateActivityResponse, error) {
	log.Printf("Received CreateActivity request: %+v", req)
	activity := convertProtoToCreateParams(req.Activity)

	createdActivity, err := h.activityService.CreateActivity(ctx, &activity)
	if err != nil {
		log.Printf("Error creating activity: %v", err)
		switch err {
		case services.ErrActivityAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, err.Error())
		case services.ErrInvalidActivityData:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to create activity")
		}
	}

	return &pb.CreateActivityResponse{
		Activity: convertModelToProto(createdActivity),
	}, nil
}

func (h *ActivityHandler) GetActivity(ctx context.Context, req *pb.GetActivityRequest) (*pb.GetActivityResponse, error) {
	activity, err := h.activityService.GetActivity(ctx, int32(req.Id))
	if err != nil {
		log.Printf("Error getting activity: %v", err)
		if err == services.ErrActivityNotFound {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, "failed to get activity")
	}

	return &pb.GetActivityResponse{
		Activity: convertModelToProto(activity),
	}, nil
}

func (h *ActivityHandler) UpdateActivity(ctx context.Context, req *pb.UpdateActivityRequest) (*pb.UpdateActivityResponse, error) {
	log.Printf("Received UpdateActivity request: %+v", req)

	// Convert Proto → sqlc params
	params := convertProtoToUpdateParams(req.Activity)

	updatedActivity, err := h.activityService.UpdateActivity(ctx, params)
	if err != nil {
		log.Printf("Error updating activity: %v", err)
		switch err {
		case services.ErrActivityNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case services.ErrInvalidActivityData:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to update activity")
		}
	}

	return &pb.UpdateActivityResponse{
		Activity: convertModelToProto(updatedActivity),
	}, nil
}

func (h *ActivityHandler) DeleteActivity(ctx context.Context, req *pb.DeleteActivityRequest) (*pb.DeleteActivityResponse, error) {
	log.Printf("Received DeleteActivity request: %+v", req)

	err := h.activityService.DeleteActivity(ctx, int32(req.Id))
	if err != nil {
		log.Printf("Error deleting activity: %v", err)
		if err == services.ErrActivityNotFound {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, "failed to delete activity")
	}

	return &pb.DeleteActivityResponse{Success: true}, nil
}

func (h *ActivityHandler) ListActivities(ctx context.Context, req *pb.ListActivitiesRequest) (*pb.ListActivitiesResponse, error) {
	activities, err := h.activityService.ListActivities(ctx, uint(req.PageNumber), uint(req.PageSize))
	if err != nil {
		log.Printf("Error listing activities: %v", err)
		if err == services.ErrInvalidActivityData {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, "failed to list activities")
	}

	var protoActivities []*pb.Activity
	for _, activity := range activities {
		protoActivities = append(protoActivities, convertModelToProto(&activity))
	}

	return &pb.ListActivitiesResponse{Activities: protoActivities}, nil
}

func convertProtoToCreateParams(proto *pb.Activity) db.CreateActivityParams {
	// Handle description
	desc := sql.NullString{Valid: proto.Description != ""}
	if desc.Valid {
		desc.String = proto.Description
	}

	// Handle due_date
	var due sql.NullTime
	if proto.DueDate != "" {
		if t, err := time.Parse(time.RFC3339, proto.DueDate); err == nil {
			due = sql.NullTime{Time: t, Valid: true}
		}
	}

	return db.CreateActivityParams{
		Title:       proto.Title,
		Description: desc,
		Type:        proto.Type,
		Status:      proto.Status,
		DueDate:     due,
		ContactID:   int32(proto.ContactId),
	}
}

func convertProtoToUpdateParams(proto *pb.Activity) db.UpdateActivityParams {
	// Handle description
	desc := sql.NullString{Valid: proto.Description != ""}
	if desc.Valid {
		desc.String = proto.Description
	}

	// Handle due_date
	var due sql.NullTime
	if proto.DueDate != "" {
		if t, err := time.Parse(time.RFC3339, proto.DueDate); err == nil {
			due = sql.NullTime{Time: t, Valid: true}
		}
	}

	return db.UpdateActivityParams{
		ID:          int32(proto.Id),
		Description: desc,
		Status:      proto.Status,
		DueDate:     due,
	}
}

// ---------- SQLC Model → Proto ----------

func convertModelToProto(model *db.Activity) *pb.Activity {
	// Handle description
	desc := ""
	if model.Description.Valid {
		desc = model.Description.String
	}

	// Handle due_date
	dueDate := ""
	if model.DueDate.Valid {
		dueDate = model.DueDate.Time.Format(time.RFC3339)
	}

	// Handle created_at / updated_at
	created := ""
	if model.CreatedAt.Valid {
		created = model.CreatedAt.Time.Format(time.RFC3339)
	}
	updated := ""
	if model.UpdatedAt.Valid {
		updated = model.UpdatedAt.Time.Format(time.RFC3339)
	}

	return &pb.Activity{
		Id:          uint32(model.ID),
		Title:       model.Title,
		Description: desc,
		Type:        model.Type,
		Status:      model.Status,
		DueDate:     dueDate,
		ContactId:   uint32(model.ContactID),
		CreatedAt:   created,
		UpdatedAt:   updated,
	}
}
