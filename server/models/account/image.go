package models

import (
	"mime/multipart"

	"github.com/Cprime50/Gopay/helper"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type ImageRepository interface {
	FileUpload(file *File) (string, error)
	RemoteUpload(url *Url) (string, error)
}

type media struct{}

func NewImageRepository() ImageRepository {
	return &media{}
}

type File struct {
	File multipart.File `json:"file,omitempty" validate:"required"`
}

type Url struct {
	Url string `json:"url,omitempty" validate:"required"`
}

func (*media) FileUpload(file *File) (string, error) {
	// Validate
	err := validate.Struct(file)
	if err != nil {
		return "", err
	}

	// Upload to Cloudinary
	uploadUrl, err := helper.ImageUploadHelper(file.File)
	if err != nil {
		return "", err
	}

	return uploadUrl, nil
}

func (*media) RemoteUpload(url *Url) (string, error) {
	// Validate
	err := validate.Struct(url)
	if err != nil {
		return "", err
	}

	// Upload to Cloudinary
	uploadUrl, err := helper.ImageUploadHelper(url.Url)
	if err != nil {
		return "", err
	}

	return uploadUrl, nil
}
