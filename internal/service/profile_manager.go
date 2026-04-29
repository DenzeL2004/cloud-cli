package service

import "mws/internal/models"

// ProfileManager defines operations for profile management.
type ProfileManager interface {
	// Create creates a new profile with the given name and data.
	// Returns an error if a profile with the same name already exists.
	Create(name string, profile *models.Profile) error

	// Get returns a profile by name.
	// Returns an error if a profile with this name is not found.
	Get(name string) (*models.Profile, error)

	// List returns names of all existing profiles.
	List() ([]string, error)

	// Delete removes a profile with the given name.
	// Returns an error if a profile with this name is not found.
	Delete(name string) error
}