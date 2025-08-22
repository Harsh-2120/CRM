package models

import (
	"time"
)

// Activity represents a high-level action or event related to customer interactions.
type Activity struct {
	Id          uint   `gorm:"primaryKey"`
	Title       string `gorm:"size:255;not null;unique"`
	Description string `gorm:"type:text"`
	Type        string `gorm:"size:50;not null"` // e.g., Call, Meeting, Email
	Status      string `gorm:"size:50;not null"` // e.g., Pending, Completed, Cancelled
	DueDate     time.Time
	CreatedAt   time.Time `grom:"autoCreateTime"`
	UpdatedAt   time.Time
	// Relationships
	ContactID uint   `gorm:"not null;index"`
	Tasks     []Task `gorm:"foreignKey:ActivityID"`
}

// Task represents a specific actionable item associated with an activity.
type Task struct {
	Id          uint   `gorm:"primaryKey"`
	Title       string `gorm:"size:255;not null;unique"`
	Description string `gorm:"type:text"`
	Status      string `gorm:"size:50;not null"` // e.g., Pending, In Progress, Completed
	Priority    string `gorm:"size:50;not null"` // e.g., Low, Medium, High
	DueDate     time.Time
	CreatedAt   time.Time `grom:"autoCreateTime"`
	UpdatedAt   time.Time
	// Relationships
	ActivityID uint     `gorm:"not null;index"`
	Activity   Activity `gorm:"foreignKey:ActivityID"`
}

type TaxationDetail struct {
	ID          uint      `gorm:"primaryKey"`
	Country     string    `gorm:"size:100;not null"` // ISO code or country name
	TaxType     string    `gorm:"size:50;not null"`  // e.g., "GST", "VAT"
	Rate        float64   `gorm:"not null"`          // e.g., 18.5 = 18.5%
	Description string    `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time
}

type Company struct {
	ID             uint      `gorm:"primaryKey"`
	Name           string    `gorm:"size:100;not null"`
	Industry       string    `gorm:"size:100"`
	Website        string    `gorm:"size:255"`
	Phone          string    `gorm:"size:20"`
	Email          string    `gorm:"size:100"`
	Address        string    `gorm:"size:255"`
	City           string    `gorm:"size:100"`
	State          string    `gorm:"size:100"`
	Country        string    `gorm:"size:100"`
	ZipCode        string    `gorm:"size:20"`
	CreatedBy      uint      // user_id
	OrganizationID uint      // org_id
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time
}

type Contact struct {
	ID          uint   `gorm:"primaryKey"`
	ContactType string `gorm:"size:20;not null"` // "individual" or "company"

	// Individual
	FirstName string `gorm:"size:100"`
	LastName  string `gorm:"size:100"`

	// Fallback when company is not registered
	CompanyName string `gorm:"size:100"`

	// Foreign key to CRM company (optional)
	CompanyID *uint `gorm:"index"`
	Company   *Company

	Email               string `gorm:"size:100;unique;not null"`
	Phone               string `gorm:"size:20"`
	Address             string `gorm:"size:255"`
	City                string `gorm:"size:100"`
	State               string `gorm:"size:100"`
	Country             string `gorm:"size:100"`
	ZipCode             string `gorm:"size:20"`
	Position            string `gorm:"size:100"`
	SocialMediaProfiles string `gorm:"type:text"`
	Notes               string `gorm:"type:text"`

	// Optional taxation details association.
	TaxationDetailID *uint
	TaxationDetail   *TaxationDetail `gorm:"foreignKey:TaxationDetailID"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time
}

type Lead struct {
	ID             uint   `gorm:"primaryKey"`
	FirstName      string `gorm:"not null"`
	LastName       string `gorm:"not null"`
	Email          string `gorm:"uniqueIndex;not null"`
	Phone          string
	Status         string    `gorm:"not null"`
	AssignedTo     int       // This will reference a user ID from user-service
	OrganizationID int       // This will reference an organization ID from organization-service
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time
}

type Opportunity struct {
	Id          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Stage       string    `json:"stage"`
	Amount      float64   `json:"amount"`
	CloseDate   time.Time `json:"close_date"`
	Probability float64   `json:"probability"`
	LeadID      uint      `json:"lead_id"`
	AccountID   uint      `json:"account_id"`
	OwnerID     uint      `json:"owner_id"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
