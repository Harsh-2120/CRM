package repository

import (
	"context"
	"crm/internal/core/domain/models"
)

type ActivityRepository interface {
	CreateActivity(ctx context.Context, activity *models.Activity) (*models.Activity, error)
	GetActivity(id uint) (*models.Activity, error)
	GetActivityByID(id uint) (*models.Activity, error)
	UpdateActivity(activity *models.Activity) (*models.Activity, error)
	DeleteActivity(id uint) error
	ListActivities(pageNumber, pageSize uint, sortBy string, ascending bool, contactID uint) ([]models.Activity, error)
}

type TaskRepository interface {
	CreateTask(task *models.Task) (*models.Task, error)
	GetTaskByID(id uint) (*models.Task, error)
	UpdateTask(task *models.Task) (*models.Task, error)
	DeleteTask(id uint) error
	ListTasks(pageNumber uint, pageSize uint, sortBy string, ascending bool, activityID uint) ([]models.Task, error)
}

type ContactRepository interface {
	Create(contact *models.Contact) (*models.Contact, error)
	GetByID(id uint) (*models.Contact, error)
	Update(contact *models.Contact) (*models.Contact, error)
	Delete(id uint) error
	List(pageNumber uint, pageSize uint, sortBy string, ascending bool) ([]models.Contact, error)
}

type CompanyRepository interface {
	Create(company *models.Company) (*models.Company, error)
	GetByID(id uint) (*models.Company, error)
	Update(company *models.Company) (*models.Company, error)
	Delete(id uint) error
	List(orgID uint, page, size uint, sortBy string, asc bool) ([]models.Company, error)
}

type LeadRepository interface {
	Create(lead *models.Lead) (*models.Lead, error)
	GetByID(id uint) (*models.Lead, error)
	GetByEmail(email string) (*models.Lead, error) // New method for finding leads by email
	Update(lead *models.Lead) (*models.Lead, error)
	Delete(id uint) error
	GetAll() ([]models.Lead, error)
}

type OpportunityRepository interface {
	Create(opportunity *models.Opportunity) (*models.Opportunity, error)
	GetByID(id uint) (*models.Opportunity, error)
	Update(opportunity *models.Opportunity) (*models.Opportunity, error)
	Delete(id uint) error
	List(ownerID uint) ([]models.Opportunity, error)
	UpdateSelective(opportunity *models.Opportunity) error
}
