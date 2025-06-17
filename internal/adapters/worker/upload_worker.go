package worker

import (
	"context"
	"fmt"

	"github.com/PatipanDev/mini-project-golang/internal/core/domain"
	"github.com/PatipanDev/mini-project-golang/internal/core/ports"
	"github.com/riverqueue/river"
)

type UploadWorker struct {
	river.WorkerDefaults[domain.UploadProfile]
	service ports.UplaodProfileService
}

func NewUploadWorker(service ports.UplaodProfileService) *UploadWorker {
	return &UploadWorker{
		service: service,
	}
}

func (w *UploadWorker) Work(ctx context.Context, job *river.Job[domain.UploadProfile]) error {
	err := w.service.UploadProfile(ctx, &job.Args)
	if err != nil {
		return fmt.Errorf("work failed: %w", err)
	}
	fmt.Println("work done for profile:", job.Args.ID)
	return nil
}

func (UploadWorker) Kind() string {
	return "upload_image"
}
