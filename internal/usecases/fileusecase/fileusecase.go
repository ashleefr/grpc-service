package fileusecase

import (
	"context"
	"grpc-file-service/internal/models/modelssvc"
)

type FileRepository interface {
	Save(ctx context.Context, file *modelssvc.File) error
	FindByID(ctx context.Context, id string) (*modelssvc.File, error)
}

type UseCase struct {
	repo FileRepository
}

func New(repo FileRepository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) UploadFile(ctx context.Context, file *modelssvc.File) error {
	return uc.repo.Save(ctx, file)
}

func (uc *UseCase) GetFile(ctx context.Context, id string) (*modelssvc.File, error) {
	return uc.repo.FindByID(ctx, id)
}
