package service

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"mws/internal/models"
)

func TestCacheProfileManager_LoadsOnlyYamlProfiles(t *testing.T) {
	baseDir := t.TempDir()
	saveDir := filepath.Join(baseDir, "profiles")

	if err := os.MkdirAll(saveDir, 0755); err != nil {
		t.Fatalf("failed to prepare save directory: %v", err)
	}

	if err := os.WriteFile(filepath.Join(saveDir, "dev.yaml"), []byte("user: demo\nproject: alpha\n"), 0644); err != nil {
		t.Fatalf("failed to write YAML profile file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(saveDir, "note.txt"), []byte("ignore me"), 0644); err != nil {
		t.Fatalf("failed to write non-profile file: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(saveDir, "nested"), 0755); err != nil {
		t.Fatalf("failed to create nested directory: %v", err)
	}

	if err := os.WriteFile(filepath.Join(saveDir, "prod.yaml"), []byte("user: bob\nproject: beta\n"), 0644); err != nil {
		t.Fatalf("failed to write second YAML profile file: %v", err)
	}

	mgr, err := NewCacheProfileManager(saveDir, nil)
	if err != nil {
		t.Fatalf("failed to create profile manager: %v", err)
	}

	if _, ok := mgr.profileNames["dev"]; !ok {
		t.Error("expected profile 'dev' to be loaded")
	}
	if _, ok := mgr.profileNames["prod"]; !ok {
		t.Error("expected profile 'prod' to be loaded")
	}
	if _, ok := mgr.profileNames["note"]; ok {
		t.Error("did not expect non-YAML file 'note.txt' to be loaded as a profile")
	}
}

func TestCacheProfileManager_CreateNewProfile(t *testing.T) {
	mgr := newTestManager(t)
	profile := &models.Profile{User: "alice", Project: "new-project"}

	if err := mgr.Create("test", profile); err != nil {
		t.Fatalf("create profile: %v", err)
	}

	if _, ok := mgr.profileNames["test"]; !ok {
		t.Fatal("expected profile name in cache after create")
	}

	content, err := os.ReadFile(filepath.Join(mgr.saveDir, "test.yaml"))
	if err != nil {
		t.Fatalf("read created profile file: %v", err)
	}
	if len(content) == 0 {
		t.Fatal("expected created profile file to contain yaml data")
	}
}

func TestCacheProfileManager_CreateExistingProfile(t *testing.T) {
	mgr := newTestManager(t)
	profile := &models.Profile{User: "alice", Project: "new-project"}

	if err := mgr.Create("test", profile); err != nil {
		t.Fatalf("create initial profile: %v", err)
	}

	err := mgr.Create("test", profile)
	if !errors.Is(err, ErrProfileAlreadyExists) {
		t.Fatalf("expected ErrProfileAlreadyExists, got: %v", err)
	}
}

func TestCacheProfileManager_GetExistingProfile(t *testing.T) {
	mgr := newTestManager(t)
	expected := &models.Profile{User: "alice", Project: "project-x"}
	if err := mgr.Create("test", expected); err != nil {
		t.Fatalf("create profile: %v", err)
	}

	got, err := mgr.Get("test")
	if err != nil {
		t.Fatalf("get profile: %v", err)
	}
	
	if got.User != expected.User || got.Project != expected.Project {
		t.Fatalf("unexpected profile data: %+v", got)
	}
}

func TestCacheProfileManager_GetMissingProfile(t *testing.T) {
	mgr := newTestManager(t)

	_, err := mgr.Get("missing")
	if !errors.Is(err, ErrProfileNotFound) {
		t.Fatalf("expected ErrProfileNotFound, got: %v", err)
	}
}

func TestCacheProfileManager_DeleteExistingProfile(t *testing.T) {
	mgr := newTestManager(t)
	if err := mgr.Create("test", &models.Profile{User: "alice", Project: "project-x"}); err != nil {
		t.Fatalf("create profile: %v", err)
	}

	if err := mgr.Delete("test"); err != nil {
		t.Fatalf("delete profile: %v", err)
	}
	if _, ok := mgr.profileNames["test"]; ok {
		t.Fatal("expected profile name removed from cache")
	}
	if _, err := os.Stat(filepath.Join(mgr.saveDir, "test.yaml")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected profile file to be removed, got err: %v", err)
	}
}

func TestCacheProfileManager_DeleteMissingProfile(t *testing.T) {
	mgr := newTestManager(t)

	err := mgr.Delete("missing")
	if !errors.Is(err, ErrProfileNotFound) {
		t.Fatalf("expected ErrProfileNotFound, got: %v", err)
	}
}

func TestCacheProfileManager_ListProfiles(t *testing.T) {
	mgr := newTestManager(t)
	if err := mgr.Create("alice", &models.Profile{User: "alice-user", Project: "alpha"}); err != nil {
		t.Fatalf("create alice profile: %v", err)
	}
	if err := mgr.Create("bob", &models.Profile{User: "bob-user", Project: "beta"}); err != nil {
		t.Fatalf("create bob profile: %v", err)
	}

	profiles, err := mgr.List()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	slices.Sort(profiles)
	expected := []string{"alice", "bob"}
	if !slices.Equal(profiles, expected) {
		t.Fatalf("unexpected list result: got %v, want %v", profiles, expected)
	}
}

func newTestManager(t *testing.T) *CacheProfileManager {
	t.Helper()

	saveDir := t.TempDir()
	mgr, err := NewCacheProfileManager(saveDir, nil)
	if err != nil {
		t.Fatalf("create manager: %v", err)
	}
	return mgr
}
