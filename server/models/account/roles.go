package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleRepository struct {
	DB *gorm.DB
}

// NewAccountRepository initializes new account repository
func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{
		DB: db,
	}
}

type Role struct {
	gorm.Model
	ID          uint   `gorm:"primary_key"`
	Name        string `gorm:"size:50;not null;uniqueIndex" json:"name"`
	Description string `gorm:"size:255;not null" json:"description"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// create role
func (db *RoleRepository) CreateRole(ctx context.Context, role *Role) error {
	if err := db.DB.WithContext(ctx).Create(&role).Error; err != nil {
		return err
	}
	return nil
}

// Get one role by id
func (db *RoleRepository) GetRoleById(ctx context.Context, id uint) (*Role, error) {
	var role *Role
	if err := db.DB.WithContext(ctx).Where("id = ?", id).First(&role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

// Gets all roles
func (db *RoleRepository) GetAllRoles() ([]*Role, error) {
	var roles []*Role
	if err := db.DB.WithContext(ctx).Find(&role).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// Update role
func (db *RoleRepository) UpdateRole(role *Role) error {
	if err := db.DB.WithContext(ctx).Save(&role).Error; err != nil {
		return err
	}
	return nil
}

// Delete role
func (db *RoleRepository) DeleteRole(id uint) error {
	if err := db.DB.WithContext(ctx).Where("id = ?", id).Delete(&role).Error; err != nil {
		return err
	}
	return nil
}

// Assign roles to account
func (db *RoleRepository) AssignRole(accountID uuid.UUID, roleID uint) error {
	if err := db.DB.WithContext(ctx).Model(&Account{}).Where("id = ?", accountID).Update("role", roleID).Error; err != nil {
		return err
	}
	return nil
}

// Get account by role
func (db *RoleRepository) GetAccntByRole(roleID uint) (*Account, error) {
	var account *Account
	if err := db.DB.WithContext(ctx).Where("role_id = ?", roleID).Find(&account).Error; err != nil {
		return err
	}
	return nil
}
