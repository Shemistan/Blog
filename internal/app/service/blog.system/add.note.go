package blog_system

import (
	"context"
	"github.com/Shemistan/Blog/internal/app/model"
)

func (b *BlogSystemService) AddNote(ctx context.Context, note *model.Note) (int64, error) {
	b.logger.Info(funcNameAddNote, "running")

	res, err := b.noteRepo.AddNote(ctx, note)
	if err != nil {
		b.logger.Error(funcNameAddNote, "failed to adding note: ", err.Error())
		return 0, err
	}

	b.logger.Info(funcNameAddNote, "finished")
	return res, nil
}

func (b *BlogSystemService) ListNotes(ctx context.Context) ([]*model.Note, error) {
	b.logger.Info(funcNameShowNotes, "running")

	res, err := b.noteRepo.ListNotes(ctx)
	if err != nil {
		b.logger.Error(funcNameShowNotes, "failed to showing notes: ", err.Error())
		return nil, err
	}

	b.logger.Info(funcNameShowNotes, "finished")
	return res, nil
}
