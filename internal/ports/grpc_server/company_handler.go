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

type CompanyHandler struct {
	companyService services.CompanyService
	pb.UnimplementedCompanyServiceServer
}

func NewCompanyHandler(service services.CompanyService) *CompanyHandler {
	return &CompanyHandler{companyService: service}
}

func (h *CompanyHandler) CreateCompany(ctx context.Context, req *pb.CreateCompanyRequest) (*pb.CreateCompanyResponse, error) {
	log.Printf("Received CreateCompany request: %+v", req)

	company := convertProtoToCompanyModel(req.Company)

	created, err := h.companyService.CreateCompany(company)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateCompanyResponse{
		Company: convertCompanyModelToProto(created),
	}, nil
}

func (h *CompanyHandler) GetCompany(ctx context.Context, req *pb.GetCompanyRequest) (*pb.GetCompanyResponse, error) {
	log.Printf("Received GetCompany request: %+v", req)

	company, err := h.companyService.GetCompany(uint(req.Id))
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pb.GetCompanyResponse{
		Company: convertCompanyModelToProto(company),
	}, nil
}

func (h *CompanyHandler) UpdateCompany(ctx context.Context, req *pb.UpdateCompanyRequest) (*pb.UpdateCompanyResponse, error) {
	log.Printf("Received UpdateCompany request: %+v", req)

	company := convertProtoToCompanyModel(req.Company)

	updated, err := h.companyService.UpdateCompany(company)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UpdateCompanyResponse{
		Company: convertCompanyModelToProto(updated),
	}, nil
}

func (h *CompanyHandler) DeleteCompany(ctx context.Context, req *pb.DeleteCompanyRequest) (*pb.DeleteCompanyResponse, error) {
	log.Printf("Received DeleteCompany request: %+v", req)

	err := h.companyService.DeleteCompany(uint(req.Id))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeleteCompanyResponse{Success: true}, nil
}

func (h *CompanyHandler) ListCompanies(ctx context.Context, req *pb.ListCompaniesRequest) (*pb.ListCompaniesResponse, error) {
	log.Printf("Received ListCompanies request: %+v", req)

	companies, err := h.companyService.ListCompanies(
		uint(req.OrganizationId),
		uint(req.PageNumber),
		uint(req.PageSize),
		req.SortBy,
		req.Ascending,
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var protoCompanies []*pb.Company
	for _, c := range companies {
		protoCompanies = append(protoCompanies, convertCompanyModelToProto(&c))
	}

	return &pb.ListCompaniesResponse{
		Companies: protoCompanies,
	}, nil
}

func convertProtoToCompanyModel(pb *pb.Company) *models.Company {
	return &models.Company{
		ID:             uint(pb.Id),
		Name:           pb.Name,
		Industry:       pb.Industry,
		Website:        pb.Website,
		Phone:          pb.Phone,
		Email:          pb.Email,
		Address:        pb.Address,
		City:           pb.City,
		State:          pb.State,
		Country:        pb.Country,
		ZipCode:        pb.ZipCode,
		CreatedBy:      uint(pb.CreatedBy),
		OrganizationID: uint(pb.OrganizationId),
		CreatedAt:      parseTime(pb.CreatedAt),
		UpdatedAt:      parseTime(pb.UpdatedAt),
	}
}

func convertCompanyModelToProto(m *models.Company) *pb.Company {
	return &pb.Company{
		Id:             uint32(m.ID),
		Name:           m.Name,
		Industry:       m.Industry,
		Website:        m.Website,
		Phone:          m.Phone,
		Email:          m.Email,
		Address:        m.Address,
		City:           m.City,
		State:          m.State,
		Country:        m.Country,
		ZipCode:        m.ZipCode,
		CreatedBy:      uint32(m.CreatedBy),
		OrganizationId: uint32(m.OrganizationID),
		CreatedAt:      m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      m.UpdatedAt.Format(time.RFC3339),
	}
}
