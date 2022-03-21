package blog_v1

import (
	"context"
	"fmt"

	"github.com/Shemistan/Blog/internal/app/model"
	pb "github.com/Shemistan/Blog/pkg/blog.v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (b *Blog) AddNoteV1(ctx context.Context, req *pb.AddNoteV1Request) (*pb.AddNoteV1Response, error) {
	res, err := b.BlogService.AddNote(ctx, &model.Note{
		Title: req.GetTitle(),
		Text:  req.GetText(),
		Tag:   req.GetTeg(),
	})

	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("failed to adding note: %v", err.Error()))
	}

	return &pb.AddNoteV1Response{Id: res.NoteId}, nil
}
