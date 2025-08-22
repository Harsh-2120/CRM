package handler

import (
	"context"
	"crm/api/proto/pb"
	"crm/internal/core/services"
	"crm/internal/middleware/websockets"

	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LeadHandler struct {
	leadService services.LeadService
	wsServer    *websockets.Server
	pb.UnimplementedLeadServiceServer
}

func NewLeadHandler(service services.LeadService, wsServer *websockets.Server) *LeadHandler {
	return &LeadHandler{leadService: service, wsServer: wsServer}
}

func (h *LeadHandler) CreateLead(ctx context.Context, req *leadspb.CreateLeadRequest) (*leadspb.CreateLeadResponse, error) {
	log.Printf("Received CreateLead request: %+v", req)

	// Convert the protobuf Lead to the internal model Lead
	lead := ConvertProtoToModelLead(req.Lead)
	if lead == nil {
		log.Println("Error: ConvertProtoToModelLead returned nil")
		return nil, status.Error(codes.InvalidArgument, "Invalid lead data")
	}

	// Call the service layer to create the lead
	createdLead, err := h.leadService.CreateLead(lead)
	if err != nil {
		log.Printf("Error creating lead: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Notify WebSocket clients
	h.wsServer.BroadcastMessage([]byte("New lead created!"))

	// Convert the internal model Lead back to the protobuf Lead
	protoLead := ConvertModelToProtoLead(createdLead)
	if protoLead == nil {
		log.Println("Error: ConvertModelToProtoLead returned nil")
		return nil, status.Error(codes.Internal, "Failed to convert lead")
	}

	response := &leadspb.CreateLeadResponse{
		Lead: protoLead,
	}

	log.Printf("Returning CreateLead response: %+v", response)
	return response, nil
}

func (h *LeadHandler) GetLead(ctx context.Context, req *leadspb.GetLeadRequest) (*leadspb.GetLeadResponse, error) {
	lead, err := h.leadService.GetLead(uint(req.Id))
	if err != nil {
		return nil, err
	}
	return &leadspb.GetLeadResponse{
		Lead: ConvertModelToProtoLead(lead),
	}, nil
}

func (h *LeadHandler) UpdateLead(ctx context.Context, req *leadspb.UpdateLeadRequest) (*leadspb.UpdateLeadResponse, error) {
	// Convert the protobuf Lead to the internal model Lead
	lead := ConvertProtoToModelLead(req.Lead)

	// Call the service layer to update the lead
	updatedLead, err := h.leadService.UpdateLead(lead)
	if err != nil {
		return nil, err
	}

	// Convert the internal model Lead back to the protobuf Lead
	return &leadspb.UpdateLeadResponse{
		Lead: ConvertModelToProtoLead(updatedLead),
	}, nil
}

func (h *LeadHandler) DeleteLead(ctx context.Context, req *leadspb.DeleteLeadRequest) (*leadspb.DeleteLeadResponse, error) {
	err := h.leadService.DeleteLead(uint(req.Id))
	if err != nil {
		return nil, err
	}
	return &leadspb.DeleteLeadResponse{Success: true}, nil
}

func (h *LeadHandler) GetAllLeads(ctx context.Context, req *leadspb.GetAllLeadsRequest) (*leadspb.GetAllLeadsResponse, error) {
	leads, err := h.leadService.GetAllLeads()
	if err != nil {
		return nil, err
	}

	// Convert the list of model Leads to protobuf Leads
	var protoLeads []*leadspb.Lead
	for _, lead := range leads {
		protoLeads = append(protoLeads, ConvertModelToProtoLead(&lead))
	}

	return &leadspb.GetAllLeadsResponse{Leads: protoLeads}, nil
}

func (h *LeadHandler) GetLeadByEmail(ctx context.Context, req *leadspb.GetLeadByEmailRequest) (*leadspb.GetLeadByEmailResponse, error) {
	lead, err := h.leadService.GetLeadByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	return &leadspb.GetLeadByEmailResponse{
		Lead: ConvertModelToProtoLead(lead),
	}, nil
}
