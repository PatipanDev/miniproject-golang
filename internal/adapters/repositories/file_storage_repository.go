package repositories

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/PatipanDev/mini-project-golang/internal/core/ports"
)

type fileStorageRepository struct {
	baseUploadDir string
	baseURL       string
}

func NewFileStorageRepository(baseUploadDir, baseURL string) ports.FileStorageRepository {
	err := os.MkdirAll(baseUploadDir, os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf("Failed to create upload directory %s: %v", baseUploadDir, err))
	}
	return &fileStorageRepository{
		baseUploadDir: baseUploadDir,
		baseURL:       baseURL,
	}
}

func (r *fileStorageRepository) SaveFile(folderPath string, filename string, fileContent []byte) (string, error) {
	dir, _ := os.Getwd()
	fmt.Println("Current working directory:", dir)
	fmt.Println("Start saving file...")
	fmt.Println("fileContent size:", len(fileContent))
	fmt.Println("baseUploadDir:", r.baseUploadDir)
	fmt.Println("folderPath:", folderPath)

	fullUploadDir := filepath.Join(r.baseUploadDir, folderPath)
	fmt.Println("Full path to upload folder:", fullUploadDir)
	//สร้างโฟลเดอร์
	err := os.MkdirAll(fullUploadDir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create upload sub-directory %s: %w", fullUploadDir, err)
	}
	//fullpath
	filePath := filepath.Join(fullUploadDir, filename)
	fmt.Println("Saving to file path:", filePath)

	err = os.WriteFile(filePath, fileContent, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to save file to disk: %w", err)
	}

	fmt.Println("File saved successfully!")

	fileURL := fmt.Sprintf("%s%s/%s", r.baseURL, folderPath, filename)
	return fileURL, nil
}
