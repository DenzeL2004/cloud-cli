package service

import (
	"cloud-cli/internal/models"
	"errors"
	"log/slog"
	"os"
	"strings"

	"go.yaml.in/yaml/v4"
)

var (
	ErrProfileNotFound      = errors.New("profile not found")
	ErrProfileAlreadyExists = errors.New("profile already exists")
)

type CacheProfileManager struct {
	saveDir      string
	profileNames map[string]struct{}

	logger *slog.Logger
}

func NewCacheProfileManager(saveDir string, logger *slog.Logger) (*CacheProfileManager, error) {
	if logger == nil {
		logger = slog.Default()
	}

	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return nil, err
	}

	manager := &CacheProfileManager{
		saveDir:      saveDir,
		profileNames: make(map[string]struct{}),
		logger:       logger,
	}

	if err := manager.LoadProfiles(); err != nil {
		return nil, err
	}

	return manager, nil
}

func (m *CacheProfileManager) LoadProfiles() error {
	m.logger.Info("Loading profiles from directory", "dir", m.saveDir)

	enteries, err := os.ReadDir(m.saveDir)
	if err != nil {
		m.logger.Error("Failed to read profiles directory", "dir", m.saveDir, "error", err)
		return err
	}

	for _, entry := range enteries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()
		if !strings.HasSuffix(fileName, ".yaml") {
			continue
		}

		profileName := strings.TrimSuffix(fileName, ".yaml")
		m.profileNames[profileName] = struct{}{}
		m.logger.Debug("Loaded profile", "name", profileName)
	}

	return nil
}

func (m *CacheProfileManager) Create(name string, profile *models.Profile) error {
	m.logger.Info("Creating profile", "name", name)

	if m.isProfileExists(name) {
		m.logger.Info("Try create existing profile", "name", name)
		return ErrProfileAlreadyExists
	}

	filename := GetFileName(m.saveDir, name)
	if err := SaveProfileToFile(filename, profile); err != nil {
		m.logger.Error("Failed to save profile to file", "name", name, "error", err)
		return err
	}

	m.profileNames[name] = struct{}{}
	return nil

}

func (m *CacheProfileManager) Get(name string) (*models.Profile, error) {
	m.logger.Info("Getting profile", "name", name)

	if !m.isProfileExists(name) {
		m.logger.Info("Try get missing profile", "name", name)
		return nil, ErrProfileNotFound
	}

	filename := GetFileName(m.saveDir, name)
	profile, err := LoadProfileFromFile(filename)
	if err != nil {
		m.logger.Error("Failed to load profile from file", "name", name, "error", err)
		return nil, err
	}

	return profile, nil

}

func (m *CacheProfileManager) List() ([]string, error) {
	m.logger.Info("Listing all profiles")
	var profileNames []string

	for name := range m.profileNames {
		profileNames = append(profileNames, name)
	}

	return profileNames, nil
}

func (m *CacheProfileManager) Delete(name string) error {
	m.logger.Info("Deleting profile", "name", name)

	if !m.isProfileExists(name) {
		m.logger.Info("Try delete missing profile", "name", name)
		return ErrProfileNotFound
	}

	filename := GetFileName(m.saveDir, name)
	if err := os.Remove(filename); err != nil {
		m.logger.Error("Failed to remove profile file", "name", name, "error", err)
		return err
	}

	delete(m.profileNames, name)
	return nil
}

func (m *CacheProfileManager) isProfileExists(name string) bool {
	_, ok := m.profileNames[name]
	return ok
}

func GetFileName(dirPath string, name string) string {
	return dirPath + "/" + name + ".yaml"
}

func SaveProfileToFile(filename string, profile *models.Profile) error {
	data, err := yaml.Marshal(profile)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

func LoadProfileFromFile(filename string) (*models.Profile, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var profile models.Profile
	if err := yaml.Unmarshal(data, &profile); err != nil {
		return nil, err
	}

	return &profile, nil
}
