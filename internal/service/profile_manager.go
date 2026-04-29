package service

import "cloud-cli/internal/models"

type ProfileManager interface {
	Create(name string, profile *models.Profile) error
	Get(name string) (*models.Profile, error)
	List() ([]string, error)
	Delete(name string) error
}