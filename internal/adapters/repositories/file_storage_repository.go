package repositories

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"

	"github.com/PatipanDev/mini-project-golang/internal/core/ports"
	"github.com/nfnt/resize"
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
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}
	fmt.Println("Current working directory:", dir)
	fmt.Println("Start saving file...")
	fmt.Println("fileContent size:", len(fileContent))
	fmt.Println("baseUploadDir:", r.baseUploadDir)
	fmt.Println("folderPath:", folderPath)

	// Resize image
	img, _, err := image.Decode(bytes.NewReader(fileContent))
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	resizedImg := resize.Resize(800, 0, img, resize.Lanczos3)

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, resizedImg, &jpeg.Options{Quality: 50})
	if err != nil {
		return "", fmt.Errorf("failed to encode image: %w", err)
	}

	compressedContent := buf.Bytes()
	fmt.Println("Compressed fileContent size:", len(compressedContent))

	fullUploadDir := filepath.Join(r.baseUploadDir, folderPath)
	fmt.Println("Full path to upload folder:", fullUploadDir)

	err = os.MkdirAll(fullUploadDir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create upload sub-directory %s: %w", fullUploadDir, err)
	}

	filePath := filepath.Join(fullUploadDir, filename)
	fmt.Println("Saving to file path:", filePath)

	err = os.WriteFile(filePath, compressedContent, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to save file to disk: %w", err)
	}

	fmt.Println("File saved successfully!")

	// fileURL := fmt.Sprintf("%s%s/%s", r.baseURL, folderPath, filename)
	fileURL := fmt.Sprint(filename)
	return fileURL, nil
}
