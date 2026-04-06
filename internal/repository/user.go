package repository

import (
	"context"
	"errors"
	"log"
	"personal_finance_dashboard/internal/models"
	"time"

	"gorm.io/gorm"
)

func GetUserByEmail(db *gorm.DB, email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user *models.User
	if err := db.WithContext(ctx).Where("email=? AND deleted_at IS NULL", email).First(&user).Error; err != nil {
		log.Printf("repository.GetUserByEmail: failed for email=%s: %v", email, err)
		return nil, err
	}
	return user, nil
}

func GetUserById(db *gorm.DB, uid uint) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user *models.User
	if err := db.WithContext(ctx).Where("id=? AND deleted_at IS NULL", uid).First(&user).Error; err != nil {
		log.Printf("repository.GetUserById: failed for id=%v: %v", uid, err)
		return nil, err
	}
	return user, nil
}

func CreateUser(db *gorm.DB, user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// check if user exists
	userCheck, err := GetUserByEmail(db, user.Email)

	if userCheck != nil && err == nil {
		log.Printf("repository.CreateUser: user already exists email=%s", user.Email)
		return errors.New("user already exists")
	}

	if err := db.WithContext(ctx).Create(user).Error; err != nil {
		log.Printf("repository.CreateUser: failed to create user email=%s: %v", user.Email, err)
		return err
	}
	return nil
}

func DeleteUser(db *gorm.DB, uid uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := GetUserById(db, uid)
	if err != nil {
		log.Printf("repository.DeleteUser: user does not exist with id=%v: %v", uid, err)
		return err
	}

	if err := db.WithContext(ctx).Delete(user).Error; err != nil {
		log.Printf("repository.DeleteUser: cannot delete user right now with id=%v: %v", uid, err)
		return errors.New("cannot delete user right now. try later")
	}
	return nil
}

func GetAllUsers(db *gorm.DB) ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var users []*models.User
	if err := db.WithContext(ctx).Where("deleted_at IS NULL").Find(&users).Error; err != nil {
		log.Printf("repository.GetAllUsers: failed to fetch users: %v", err)
		return nil, err
	}
	return users, nil
}
