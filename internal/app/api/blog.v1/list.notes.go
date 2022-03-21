package blog_v1

import (
	"context"
	"fmt"

	pb "github.com/Shemistan/Blog/pkg/blog.v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (b *Blog) ListNotesV1(ctx context.Context, _ *emptypb.Empty) (*pb.ListNotesV1Response, error) {
	res, err := b.BlogService.ShowNotes(ctx)

	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("failed to showing notes: %v", err.Error()))
	}

	var notes []*pb.ListNotesV1Response_Note

	for _, note := range res {
		notes = append(notes, &pb.ListNotesV1Response_Note{
			Id:           note.Id,
			Title:        note.Title,
			Text:         note.Text,
			Tag:          note.Tag,
			CreatingData: note.CreatingData,
		})
	}

	return &pb.ListNotesV1Response{
		Notes: notes,
	}, nil
}
