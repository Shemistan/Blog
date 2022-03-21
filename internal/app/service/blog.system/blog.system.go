package blog_system

import (
	"context"

	"github.com/Shemistan/Blog/internal/app/config"
	"github.com/Shemistan/Blog/internal/app/model"
	"github.com/Shemistan/Blog/internal/app/repo"
	"github.com/Shemistan/Blog/internal/app/service/logger"
)

const (
	funcNameAddNote   = "AddNote"
	funcNameShowNotes = "ShowNotes"
)

type IBlogSystemService interface {
	AddNote(ctx context.Context, note *model.Note) (*model.AddResponse, error)
	ShowNotes(ctx context.Context) ([]*model.Note, error)
}

type BlogSystemService struct {
	logger    *logger.Service
	appConfig *config.Config
	noteRepo  repo.Repo
}

func NewBlogSystemService(
	logger *logger.Service,
	appConfig *config.Config,
	noteRepo repo.Repo) IBlogSystemService {

	return &BlogSystemService{
		logger:    logger,
		appConfig: appConfig,
		noteRepo:  noteRepo,
	}
}
