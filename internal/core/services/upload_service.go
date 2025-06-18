package services

import (
	"context"

	"github.com/PatipanDev/mini-project-golang/internal/core/domain"
	"github.com/PatipanDev/mini-project-golang/internal/core/ports"
)

type uploadService struct {
	repo ports.UploadProfileRepository
}

func NewUploadService(repo ports.UploadProfileRepository) ports.UplaodProfileService {
	return &uploadService{repo: repo}
}

func (s *uploadService) UploadProfile(ctx context.Context, profile *domain.UploadProfile) error {
	return s.repo.Save(ctx, profile)
}
