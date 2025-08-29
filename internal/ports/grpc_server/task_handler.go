package handler

import (
	"context"
	"crm/api/proto/pb"

	"crm/internal/core/services"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TaskHandler struct {
	taskService services.TaskService
	pb.UnimplementedTaskServiceServer
}

func NewTaskHandler(service services.TaskService) *TaskHandler {
	return &TaskHandler{taskService: service}
}

// CreateTask handles the creation of a new task.
func (h *TaskHandler) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	log.Printf("Received CreateTask request: %+v", req)

	// Convert Proto to Model
	task := convertProtoToModelTask(req.Task)

	// Create Task via Service
	createdTask, err := h.taskService.CreateTask(task)
	if err != nil {
		log.Printf("Error creating task: %v", err)
		switch err {
		case services.ErrTaskExists:
			return nil, status.Error(codes.AlreadyExists, err.Error())
		case services.ErrInvalidTaskData:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case services.ErrActivityNotFound:
			return nil, status.Error(codes.NotFound, "associated activity not found")
		default:
			return nil, status.Error(codes.Internal, "failed to create task")
		}
	}

	// Convert Model to Proto
	return &pb.CreateTaskResponse{
		Task: convertTaskModelToProto(createdTask),
	}, nil
}

// GetTask handles retrieval of a task by ID.
func (h *TaskHandler) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.GetTaskResponse, error) {
	task, err := h.taskService.GetTask(uint(req.Id))
	if err != nil {
		log.Printf("Error getting task: %v", err)
		switch err {
		case services.ErrTaskNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to get task")
		}
	}

	return &pb.GetTaskResponse{
		Task: convertTaskModelToProto(task),
	}, nil
}

// UpdateTask handles updating an existing task.
func (h *TaskHandler) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.UpdateTaskResponse, error) {
	log.Printf("Received UpdateTask request: %+v", req)

	// Convert Proto to Model
	task := convertProtoToModelTask(req.Task)

	// Update Task via Service
	updatedTask, err := h.taskService.UpdateTask(task)
	if err != nil {
		log.Printf("Error updating task: %v", err)
		switch err {
		case services.ErrTaskNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case services.ErrTaskExists:
			return nil, status.Error(codes.AlreadyExists, err.Error())
		case services.ErrInvalidTaskData:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to update task")
		}
	}

	// Convert Model to Proto
	return &pb.UpdateTaskResponse{
		Task: convertTaskModelToProto(updatedTask),
	}, nil
}

// DeleteTask handles deletion of a task by ID.
func (h *TaskHandler) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteTaskResponse, error) {
	log.Printf("Received DeleteTask request: %+v", req)

	// Delete Task via Service
	err := h.taskService.DeleteTask(uint(req.Id))
	if err != nil {
		log.Printf("Error deleting task: %v", err)
		switch err {
		case services.ErrTaskNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to delete task")
		}
	}

	// Return Success
	return &pb.DeleteTaskResponse{
		Success: true,
	}, nil
}

// ListTasks handles listing tasks with pagination and optional filtering.
func (h *TaskHandler) ListTasks(ctx context.Context, req *pb.ListTasksRequest) (*pb.ListTasksResponse, error) {
	tasks, err := h.taskService.ListTasks(uint(req.PageNumber), uint(req.PageSize), req.SortBy, req.Ascending, uint(req.ActivityId))
	if err != nil {
		log.Printf("Error listing tasks: %v", err)
		switch err {
		case services.ErrInvalidTaskData:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to list tasks")
		}
	}

	// Convert Models to Proto
	var protoTasks []*pb.Task
	for _, task := range tasks {
		protoTasks = append(protoTasks, convertTaskModelToProto(&task))
	}

	return &pb.ListTasksResponse{
		Tasks: protoTasks,
	}, nil
}

// Conversion Functions


func convertProtoToModelTask(protoTask *pb.Task) *models.Task {
	dueDate, _ := time.Parse(time.RFC3339, protoTask.DueDate)
	return &models.Task{
		Id:          uint(protoTask.Id),
		Title:       protoTask.Title,
		Description: protoTask.Description,
		Status:      protoTask.Status,
		Priority:    protoTask.Priority,
		DueDate:     dueDate,
		ActivityID:  uint(protoTask.ActivityId),
	}
}

func convertTaskModelToProto(modelTask *models.Task) *pb.Task {
	dueDate := ""
	if !modelTask.DueDate.IsZero() {
		dueDate = modelTask.DueDate.Format(time.RFC3339)
	}
	return &pb.Task{
		Id:          uint32(modelTask.Id),
		Title:       modelTask.Title,
		Description: modelTask.Description,
		Status:      modelTask.Status,
		Priority:    modelTask.Priority,
		DueDate:     dueDate,
		CreatedAt:   modelTask.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   modelTask.UpdatedAt.Format(time.RFC3339),
		ActivityId:  uint32(modelTask.ActivityID),
	}
}
