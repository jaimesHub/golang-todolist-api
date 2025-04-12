package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email     string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	FirstName string         `gorm:"type:varchar(100)" json:"first_name"`
	LastName  string         `gorm:"type:varchar(100)" json:"last_name"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Tasks      []Task     `gorm:"foreignKey:UserID" json:"-"`
	Activities []Activity `gorm:"foreignKey:UserID" json:"-"`
}

// Task represents a todo task
type Task struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Title       string         `gorm:"type:varchar(255);not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Status      string         `gorm:"type:varchar(50);default:'pending'" json:"status"` // pending, in_progress, completed
	Priority    int            `gorm:"default:0" json:"priority"`                        // 0: low, 1: medium, 2: high
	DueDate     *time.Time     `json:"due_date"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// Activity represents a user activity log
type Activity struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Action    string    `gorm:"type:varchar(100);not null" json:"action"`
	Entity    string    `gorm:"type:varchar(100);not null" json:"entity"`
	EntityID  uuid.UUID `gorm:"type:uuid" json:"entity_id"`
	Details   string    `gorm:"type:text" json:"details"`
	CreatedAt time.Time `json:"created_at"`
}

// BeforeCreate is a GORM hook that runs before creating a record
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// BeforeCreate is a GORM hook that runs before creating a record
func (t *Task) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// BeforeCreate is a GORM hook that runs before creating a record
func (a *Activity) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
