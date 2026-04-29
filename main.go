package main

import (
	"fmt"
	"log/slog"
	"mws/cmd"
	"mws/internal/service"
	"os"
	"path/filepath"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed:", err)
		return
	}

	logFilePath := filepath.Join(os.TempDir(), "mws", "mws.log")
	if err := os.MkdirAll(filepath.Dir(logFilePath), 0o755); err != nil {
		fmt.Println("Failed:", err)
		return
	}

	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Println("Failed:", err)
		return
	}
	defer logFile.Close()

	logger := slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{Level: slog.LevelDebug}))

	profManager, err := service.NewCacheProfileManager(dir, logger)
	if err != nil {
		fmt.Println("Failed:", err)
		return
	}

	cmd.Execute(profManager)
}
