package handler

import (
	"context"
	"crm/api/proto/pb"
	"crm/internal/core/domain/models"
	"crm/internal/core/services"
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

// CreateActivity handles the creation of a new activity.
func (h *ActivityHandler) CreateActivity(ctx context.Context, req *pb.CreateActivityRequest) (*pb.CreateActivityResponse, error) {
	log.Printf("Received CreateActivity request: %+v", req)

	// Convert Proto to Model
	activity := convertProtoToModel(req.Activity)

	// Create Activity via Service
	createdActivity, err := h.activityService.CreateActivity(ctx, activity)
	if err != nil {
		log.Printf("Error creating activity: %v", err)
		switch err {
		case services.ErrActivityExists:
			return nil, status.Error(codes.AlreadyExists, err.Error())
		case services.ErrInvalidActivityData:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to create activity")
		}
	}

	// Convert Model to Proto
	return &pb.CreateActivityResponse{
		Activity: convertModelToProto(createdActivity),
	}, nil
}

// GetActivity handles retrieval of an activity by ID.
func (h *ActivityHandler) GetActivity(ctx context.Context, req *pb.GetActivityRequest) (*pb.GetActivityResponse, error) {
	activity, err := h.activityService.GetActivity(uint(req.Id))
	if err != nil {
		log.Printf("Error getting activity: %v", err)
		switch err {
		case services.ErrActivityNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to get activity")
		}
	}

	return &pb.GetActivityResponse{
		Activity: convertModelToProto(activity),
	}, nil
}

// UpdateActivity handles updating an existing activity.
func (h *ActivityHandler) UpdateActivity(ctx context.Context, req *pb.UpdateActivityRequest) (*pb.UpdateActivityResponse, error) {

	log.Printf("Received UpdateActivity request: %+v", req)

	// Convert Proto to Model
	activity := convertProtoToModel(req.Activity)

	// Update Activity via Service
	updatedActivity, err := h.activityService.UpdateActivity(activity)

	if err != nil {
		log.Printf("Error updating activity: %v", err)
		switch err {
		case services.ErrActivityNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case services.ErrActivityExists:
			return nil, status.Error(codes.AlreadyExists, err.Error())
		case services.ErrInvalidActivityData:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to update activity")
		}
	}

	// Convert Model to Proto
	log.Print("updated log", updatedActivity)
	return &pb.UpdateActivityResponse{
		Activity: convertModelToProto(updatedActivity),
	}, nil
}

// DeleteActivity handles deletion of an activity by ID.
func (h *ActivityHandler) DeleteActivity(ctx context.Context, req *pb.DeleteActivityRequest) (*pb.DeleteActivityResponse, error) {
	log.Printf("Received DeleteActivity request: %+v", req)

	// Delete Activity via Service
	err := h.activityService.DeleteActivity(uint(req.Id))
	if err != nil {
		log.Printf("Error deleting activity: %v", err)
		switch err {
		case services.ErrActivityNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to delete activity")
		}
	}

	// Return Success
	return &pb.DeleteActivityResponse{
		Success: true,
	}, nil
}

// ListActivities handles listing activities with pagination and optional filtering.
func (h *ActivityHandler) ListActivities(ctx context.Context, req *pb.ListActivitiesRequest) (*pb.ListActivitiesResponse, error) {
	activities, err := h.activityService.ListActivities(uint(req.PageNumber), uint(req.PageSize), req.SortBy, req.Ascending, uint(req.ContactId))
	if err != nil {
		log.Printf("Error listing activities: %v", err)
		switch err {
		case services.ErrInvalidActivityData:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to list activities")
		}
	}

	// Convert Models to Proto
	var protoActivities []*pb.Activity
	for _, activity := range activities {
		protoActivities = append(protoActivities, convertModelToProto(&activity))
	}

	return &pb.ListActivitiesResponse{
		Activities: protoActivities,
	}, nil
}

// Conversion Functions

func convertProtoToModel(protoActivity *pb.Activity) *models.Activity {
	dueDate, _ := time.Parse(time.RFC3339, protoActivity.DueDate)
	return &models.Activity{
		Id:          uint(protoActivity.Id),
		Title:       protoActivity.Title,
		Description: protoActivity.Description,
		Type:        protoActivity.Type,
		Status:      protoActivity.Status,
		DueDate:     dueDate,
		ContactID:   uint(protoActivity.ContactId),
	}
}

func convertModelToProto(modelActivity *models.Activity) *pb.Activity {
	dueDate := ""
	if !modelActivity.DueDate.IsZero() {
		dueDate = modelActivity.DueDate.Format(time.RFC3339)
	}
	return &pb.Activity{
		Id:          uint32(modelActivity.Id),
		Title:       modelActivity.Title,
		Description: modelActivity.Description,
		Type:        modelActivity.Type,
		Status:      modelActivity.Status,
		DueDate:     dueDate,
		CreatedAt:   modelActivity.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   modelActivity.UpdatedAt.Format(time.RFC3339),
		ContactId:   uint32(modelActivity.ContactID),
	}
}
