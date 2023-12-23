package models

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Cprime50/Gopay/helper"
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

// Get one role by id
func (db *RoleRepository) GetRoleById(ctx context.Context, id uint) (*Role, error) {
	var role *Role
	if err := db.DB.WithContext(ctx).Where("id = ?", id).First(&role).Error; err != nil {
		if err != nil {
			if errors.Is(gorm.ErrRecordNotFound, err) {
				fmt.Printf("role with ID %d not found\n", role.ID)
				return nil, helper.NewNotFound("id", fmt.Sprintf("%d", role.ID))
			}
		}
		log.Println("Error quering db", err)
		return nil, helper.NewInternal()
	}
	return role, nil
}

// Gets all roles
func (db *RoleRepository) GetAllRoles(ctx context.Context) ([]*Role, error) {
	var roles []*Role
	if err := db.DB.WithContext(ctx).Find(&roles).Error; err != nil {
		log.Println("Error quering db", err)
		return nil, helper.NewInternal()
	}
	return roles, nil
}

// Assign roles to account
func (db *RoleRepository) AssignRole(ctx context.Context, accountID uuid.UUID, roleID uint) error {
	if err := db.DB.WithContext(ctx).Model(&Account{}).Where("id = ?", accountID).Update("role", roleID).Error; err != nil {
		if err != nil {
			if errors.Is(gorm.ErrRecordNotFound, err) {
				fmt.Printf("account with ID %s not found\n", accountID)
				return helper.NewNotFound("id", accountID.String())
			}
		}
		log.Println("Error quering db", err)
		return helper.NewInternal()
	}
	return nil
}

// Get accounts by role
func (db *RoleRepository) GetAccntByRole(ctx context.Context, roleID uint) (*[]Account, error) {
	var accounts *[]Account
	if err := db.DB.WithContext(ctx).Where("role_id = ?", roleID).Find(&accounts).Error; err != nil {
		if err != nil {
			if errors.Is(gorm.ErrRecordNotFound, err) {
				fmt.Printf("role with ID %d not found\n", roleID)
				return nil, helper.NewNotFound("id", fmt.Sprintf("%d", roleID))
			}
		}
		log.Println("Error quering db", err)
		return nil, helper.NewInternal()
	}
	return accounts, nil
}
