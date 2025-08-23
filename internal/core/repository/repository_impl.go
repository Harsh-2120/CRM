// package repository

// import (
// 	"context"
// 	"crm/internal/core/domain/models"
// 	"errors"
// 	"fmt"
// 	"log"
// 	"time"

// 	"github.com/jackc/pgconn"
// 	"gorm.io/gorm"
// )

// // ====================== ERRORS ======================
// var (
// 	// Activity
// 	ErrActivityExists   = errors.New("activity with this title already exists")
// 	ErrActivityNotFound = errors.New("activity not found")
// 	ErrTaskExists       = errors.New("task with this title already exists")
// 	ErrTaskNotFound     = errors.New("task not found")

// 	// Company
// 	ErrCompanyNotFound = errors.New("company not found")

// 	// Contact
// 	ErrContactExists   = errors.New("contact with this email already exists")
// 	ErrContactNotFound = errors.New("contact not found")
// )

// // ====================== ACTIVITY REPOSITORY ======================
// type baseRepository struct {
// 	db *gorm.DB
// }
// type companyRepository struct {
// 	*baseRepository
// }

// type activityRepository struct {
// 	*baseRepository
// }

// type taskRepository struct {
// 	*baseRepository
// }

// type contactRepository struct {
// 	*baseRepository
// }

// type leadRepository struct {
// 	*baseRepository
// }

// type opportunityRepository struct {
// 	*baseRepository
// }

// func newBaseRepository(db *gorm.DB) *baseRepository {
// 	return &baseRepository{db: db}
// }

// func NewCompanyRepository(db *gorm.DB) CompanyRepository {
// 	return &companyRepository{newBaseRepository(db)}
// }

// func NewActivityRepository(db *gorm.DB) ActivityRepository {
// 	return &activityRepository{newBaseRepository(db)}
// }

// func NewTaskRepository(db *gorm.DB) TaskRepository {
// 	return &taskRepository{newBaseRepository(db)}
// }

// func NewContactRepository(db *gorm.DB) ContactRepository {
// 	return &contactRepository{newBaseRepository(db)}
// }

// func NewLeadRepository(db *gorm.DB) LeadRepository {
// 	return &leadRepository{newBaseRepository(db)}
// }

// func NewOpportunityRepository(db *gorm.DB) OpportunityRepository {
// 	return &opportunityRepository{newBaseRepository(db)}
// }

// // ActivityRepository methods

// func (r *activityRepository) CreateActivity(ctx context.Context, activity *models.Activity) (*models.Activity, error) {
// 	tx := r.db.WithContext(ctx).Begin() // Start transaction

// 	// Ensure rollback happens if function exits unexpectedly
// 	defer func() {
// 		if r := recover(); r != nil {
// 			tx.Rollback()
// 			log.Printf("Transaction panic recovered: %v", r)
// 		}
// 	}()

// 	// Attempt to create activity
// 	if err := tx.Create(activity).Error; err != nil {
// 		tx.Rollback()                                                        // Explicit rollback on failure
// 		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" { // Unique violation
// 			return nil, ErrActivityExists
// 		}
// 		log.Printf("Failed to create activity: %v", err)
// 		return nil, err
// 	}

// 	// Commit transaction on success
// 	if err := tx.Commit().Error; err != nil {
// 		log.Printf("Transaction commit failed: %v", err)
// 		return nil, err
// 	}

// 	return activity, nil
// }

// func (r *activityRepository) GetActivity(id uint) (*models.Activity, error) {
// 	var activity models.Activity
// 	if err := r.db.Preload("Tasks").First(&activity, id).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, ErrActivityNotFound
// 		}
// 		return nil, err
// 	}
// 	return &activity, nil
// }

// // UpdateActivity modifies an existing activity.
// func (r *activityRepository) UpdateActivity(activity *models.Activity) (*models.Activity, error) {
// 	// Perform the update while omitting CreatedAt
// 	result := r.db.Model(&models.Activity{}).Where("id = ?", activity.Id).Omit("CreatedAt").Updates(activity)

// 	// Check for errors
// 	if result.Error != nil {
// 		if isUniqueConstraintError(result.Error, "activities_title_key") {
// 			return nil, ErrActivityExists
// 		}
// 		return nil, result.Error
// 	}

// 	// Check if any row was affected
// 	if result.RowsAffected == 0 {
// 		return nil, ErrActivityNotFound
// 	}

// 	// Reload the updated record
// 	var updatedActivity models.Activity
// 	if err := r.db.First(&updatedActivity, "id = ?", activity.Id).Error; err != nil {
// 		return nil, err
// 	}

// 	return &updatedActivity, nil
// }

// // DeleteActivity removes an activity by its Id.
// func (r *activityRepository) DeleteActivity(id uint) error {
// 	result := r.db.Delete(&models.Activity{}, id)
// 	if result.Error != nil {
// 		return result.Error
// 	}
// 	if result.RowsAffected == 0 {
// 		return ErrActivityNotFound
// 	}
// 	return nil
// }

// // ListActivities retrieves activities with pagination, sorting, and optional filtering by contact.
// func (r *activityRepository) ListActivities(pageNumber uint, pageSize uint, sortBy string, ascending bool, contactID uint) ([]models.Activity, error) {
// 	var activities []models.Activity

// 	query := r.db.Model(&models.Activity{})

// 	// Apply filter by contact Id if provided
// 	if contactID != 0 {
// 		query = query.Where("contact_id = ?", contactID)
// 	}

// 	// Apply sorting
// 	if sortBy != "" {
// 		order := sortBy
// 		if ascending {
// 			order += " ASC"
// 		} else {
// 			order += " DESC"
// 		}
// 		query = query.Order(order)
// 	}

// 	// Apply pagination
// 	offset := (pageNumber - 1) * pageSize
// 	query = query.Offset(int(offset)).Limit(int(pageSize))

// 	if err := query.Preload("Tasks").Find(&activities).Error; err != nil {
// 		return nil, err
// 	}

// 	return activities, nil
// }

// // taskRepository methods

// func (r *activityRepository) GetActivityByID(id uint) (*models.Activity, error) {
// 	var activity models.Activity
// 	if err := r.db.First(&activity, id).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, ErrActivityNotFound
// 		}
// 		return nil, err
// 	}
// 	return &activity, nil
// }

// // CreateTask inserts a new task into the database.
// func (r *taskRepository) CreateTask(task *models.Task) (*models.Task, error) {
// 	if err := r.db.Create(task).Error; err != nil {
// 		if isUniqueConstraintError(err, "tasks_title_key") {
// 			return nil, ErrTaskExists
// 		}
// 		return nil, err
// 	}
// 	return task, nil
// }

// // GetTaskByID retrieves a task by its Id.
// func (r *taskRepository) GetTaskByID(id uint) (*models.Task, error) {
// 	var task models.Task
// 	if err := r.db.First(&task, id).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, ErrTaskNotFound
// 		}
// 		return nil, err
// 	}
// 	return &task, nil
// }

// // UpdateTask modifies an existing task.
// func (r *taskRepository) UpdateTask(task *models.Task) (*models.Task, error) {
// 	// Perform the update, omitting CreatedAt
// 	result := r.db.Model(&models.Task{}).Where("id = ?", task.Id).Omit("CreatedAt").Updates(task)

// 	// Check for update errors
// 	if result.Error != nil {
// 		if isUniqueConstraintError(result.Error, "tasks_title_key") {
// 			return nil, ErrTaskExists
// 		}
// 		return nil, result.Error
// 	}

// 	// Check if any row was actually updated
// 	if result.RowsAffected == 0 {
// 		return nil, ErrTaskNotFound
// 	}

// 	// Reload the updated record
// 	var updatedTask models.Task
// 	if err := r.db.First(&updatedTask, "id = ?", task.Id).Error; err != nil {
// 		log.Println("Failed to reload updated task:", err)
// 		return nil, err
// 	}

// 	return &updatedTask, nil
// }

// // DeleteTask removes a task by its Id.
// func (r *taskRepository) DeleteTask(id uint) error {
// 	result := r.db.Delete(&models.Task{}, id)
// 	if result.Error != nil {
// 		return result.Error
// 	}
// 	if result.RowsAffected == 0 {
// 		return ErrTaskNotFound
// 	}
// 	return nil
// }

// // ListTasks retrieves tasks with pagination, sorting, and optional filtering by activity.
// func (r *taskRepository) ListTasks(pageNumber uint, pageSize uint, sortBy string, ascending bool, activityID uint) ([]models.Task, error) {
// 	var tasks []models.Task

// 	query := r.db.Model(&models.Task{})

// 	// Apply filter by activity Id if provided
// 	if activityID != 0 {
// 		query = query.Where("activity_id = ?", activityID)
// 	}

// 	// Apply sorting
// 	if sortBy != "" {
// 		order := sortBy
// 		if ascending {
// 			order += " ASC"
// 		} else {
// 			order += " DESC"
// 		}
// 		query = query.Order(order)
// 	}

// 	// Apply pagination
// 	offset := (pageNumber - 1) * pageSize
// 	query = query.Offset(int(offset)).Limit(int(pageSize))

// 	if err := query.Find(&tasks).Error; err != nil {
// 		return nil, err
// 	}

// 	return tasks, nil
// }

// // companyRepository methods

// func (r *companyRepository) Create(company *models.Company) (*models.Company, error) {
// 	if err := r.db.Create(company).Error; err != nil {
// 		return nil, err
// 	}
// 	return company, nil
// }

// func (r *companyRepository) GetByID(id uint) (*models.Company, error) {
// 	var company models.Company
// 	if err := r.db.First(&company, id).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, ErrCompanyNotFound
// 		}
// 		return nil, err
// 	}
// 	return &company, nil
// }

// func (r *companyRepository) Update(company *models.Company) (*models.Company, error) {
// 	err := r.db.Model(&models.Company{}).Where("id = ?", company.ID).Updates(company).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return company, nil
// }

// func (r *companyRepository) Delete(id uint) error {
// 	result := r.db.Delete(&models.Company{}, id)
// 	if result.Error != nil {
// 		return result.Error
// 	}
// 	if result.RowsAffected == 0 {
// 		return ErrCompanyNotFound
// 	}
// 	return nil
// }

// func (r *companyRepository) List(orgID uint, page, size uint, sortBy string, asc bool) ([]models.Company, error) {
// 	var companies []models.Company

// 	offset := (page - 1) * size
// 	order := sortBy
// 	if sortBy != "" {
// 		if asc {
// 			order += " ASC"
// 		} else {
// 			order += " DESC"
// 		}
// 	} else {
// 		order = "created_at DESC"
// 	}

// 	err := r.db.
// 		Where("organization_id = ?", orgID).
// 		Order(order).
// 		Offset(int(offset)).
// 		Limit(int(size)).
// 		Find(&companies).Error

// 	if err != nil {
// 		return nil, err
// 	}

// 	return companies, nil
// }

// // contactRepository methods

// func (r *contactRepository) Create(contact *models.Contact) (*models.Contact, error) {
// 	if err := r.db.Create(contact).Error; err != nil {
// 		if isUniqueConstraintError(err, "contacts_email_key") {
// 			return nil, ErrContactExists
// 		}
// 		return nil, err
// 	}
// 	return contact, nil
// }

// // GetByID retrieves a unified contact by its ID.
// func (r *contactRepository) GetByID(id uint) (*models.Contact, error) {
// 	var contact models.Contact
// 	if err := r.db.First(&contact, id).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, ErrContactNotFound
// 		}
// 		return nil, err
// 	}
// 	return &contact, nil
// }

// // Update modifies an existing unified contact.
// func (r *contactRepository) Update(contact *models.Contact) (*models.Contact, error) {
// 	result := r.db.Model(&models.Contact{}).Where("id = ?", contact.ID).Updates(contact)
// 	if result.Error != nil {
// 		if isUniqueConstraintError(result.Error, "contacts_email_key") {
// 			return nil, ErrContactExists
// 		}
// 		return nil, result.Error
// 	}
// 	if result.RowsAffected == 0 {
// 		return nil, ErrContactNotFound
// 	}
// 	var updatedContact models.Contact
// 	if err := r.db.First(&updatedContact, "id=?", contact.ID).Error; err != nil {
// 		log.Printf("failed to reload updated contact")
// 		return nil, err
// 	}
// 	return &updatedContact, nil
// }

// // Delete removes a unified contact by its ID.
// func (r *contactRepository) Delete(id uint) error {
// 	result := r.db.Delete(&models.Contact{}, id)
// 	if result.Error != nil {
// 		return result.Error
// 	}
// 	if result.RowsAffected == 0 {
// 		return ErrContactNotFound
// 	}
// 	return nil
// }

// // List retrieves contacts with pagination and sorting.
// // This function works with the unified contact model that may include additional fields
// // like ContactType, CompanyName, and TaxationDetailID.
// func (r *contactRepository) List(pageNumber uint, pageSize uint, sortBy string, ascending bool) ([]models.Contact, error) {
// 	var contacts []models.Contact

// 	query := r.db.Model(&models.Contact{})

// 	// Apply sorting if provided.
// 	if sortBy != "" {
// 		order := sortBy
// 		if ascending {
// 			order += " ASC"
// 		} else {
// 			order += " DESC"
// 		}
// 		query = query.Order(order)
// 	}

// 	// Apply pagination.
// 	offset := (pageNumber - 1) * pageSize
// 	query = query.Offset(int(offset)).Limit(int(pageSize))

// 	if err := query.Find(&contacts).Error; err != nil {
// 		return nil, err
// 	}

// 	return contacts, nil
// }

// // isUniqueConstraintError checks for PostgreSQL unique constraint violations.
// func isUniqueConstraintError(err error, constraintName string) bool {
// 	var pgErr *pgconn.PgError
// 	if errors.As(err, &pgErr) {
// 		if pgErr.Code == "23505" && pgErr.ConstraintName == constraintName {
// 			return true
// 		}
// 	}
// 	return false
// }

// // leadRepository methods

// func (r *leadRepository) Create(lead *models.Lead) (*models.Lead, error) {
// 	if err := r.db.Create(lead).Error; err != nil {
// 		return nil, err
// 	}
// 	fmt.Println("Lead created successfully")
// 	return lead, nil
// }

// func (r *leadRepository) GetByID(id uint) (*models.Lead, error) {
// 	var lead models.Lead
// 	if err := r.db.First(&lead, id).Error; err != nil {
// 		return nil, err
// 	}
// 	return &lead, nil
// }

// // New function to get lead by email
// func (r *leadRepository) GetByEmail(email string) (*models.Lead, error) {
// 	var lead models.Lead
// 	if err := r.db.Where("email = ?", email).First(&lead).Error; err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}
// 	return &lead, nil
// }

// func (r *leadRepository) Update(lead *models.Lead) (*models.Lead, error) {
// 	// if err := r.db.Save(lead).Error; err != nil {
// 	// 	return nil, err
// 	// }
// 	// return lead, nil

// 	if err := r.db.Where("id=?", lead.ID).Updates(lead).Error; err != nil {
// 		print("failed to update")
// 		return nil, err
// 	}

// 	var updatedLead models.Lead

// 	if err := r.db.First(&updatedLead, "id=?", lead.ID).Error; err != nil {
// 		print("failed to reload the updated lead")
// 		return lead, nil
// 	}

// 	return &updatedLead, nil
// }

// func (r *leadRepository) Delete(id uint) error {
// 	if err := r.db.Delete(&models.Lead{}, id).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (r *leadRepository) GetAll() ([]models.Lead, error) {
// 	var leads []models.Lead
// 	if err := r.db.Find(&leads).Error; err != nil {
// 		return nil, err
// 	}
// 	return leads, nil
// }

// // opportunityRepository methods

// func (r *opportunityRepository) Create(opportunity *models.Opportunity) (*models.Opportunity, error) {
// 	if err := r.db.Create(opportunity).Error; err != nil {
// 		return nil, err
// 	}
// 	return opportunity, nil
// }

// func (r *opportunityRepository) GetByID(id uint) (*models.Opportunity, error) {
// 	var opportunity models.Opportunity
// 	if err := r.db.First(&opportunity, id).Error; err != nil {
// 		return nil, err
// 	}
// 	return &opportunity, nil
// }

// func (r *opportunityRepository) Update(opportunity *models.Opportunity) (*models.Opportunity, error) {
// 	if err := r.db.Save(opportunity).Error; err != nil {
// 		return nil, err
// 	}
// 	return opportunity, nil
// }

// func (r *opportunityRepository) Delete(id uint) error {
// 	if err := r.db.Delete(&models.Opportunity{}, id).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (r *opportunityRepository) List(ownerID uint) ([]models.Opportunity, error) {
// 	var opportunities []models.Opportunity
// 	query := r.db
// 	if ownerID != 0 {
// 		query = query.Where("owner_id = ?", ownerID)
// 	}
// 	if err := query.Find(&opportunities).Error; err != nil {
// 		return nil, err
// 	}
// 	return opportunities, nil
// }

// func (r *opportunityRepository) UpdateSelective(opportunity *models.Opportunity) error {
// 	// Create a map of fields to update
// 	updates := map[string]interface{}{}

// 	if opportunity.Name != "" {
// 		updates["name"] = opportunity.Name
// 	}
// 	if opportunity.Description != "" {
// 		updates["description"] = opportunity.Description
// 	}
// 	if opportunity.Stage != "" {
// 		updates["stage"] = opportunity.Stage
// 	}
// 	if opportunity.Amount != 0 {
// 		updates["amount"] = opportunity.Amount
// 	}
// 	if !opportunity.CloseDate.IsZero() {
// 		updates["close_date"] = opportunity.CloseDate
// 	}
// 	if opportunity.Probability != 0 {
// 		updates["probability"] = opportunity.Probability
// 	}
// 	if opportunity.LeadID != 0 {
// 		updates["lead_id"] = opportunity.LeadID
// 	}
// 	if opportunity.AccountID != 0 {
// 		updates["account_id"] = opportunity.AccountID
// 	}
// 	if opportunity.OwnerID != 0 {
// 		updates["owner_id"] = opportunity.OwnerID
// 	}

// 	updates["updated_at"] = time.Now()

// 	return r.db.Model(&models.Opportunity{}).Where("id = ?", opportunity.Id).Updates(updates).Error
// }
