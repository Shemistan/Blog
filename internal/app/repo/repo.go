package repo

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/Shemistan/Blog/internal/app/model"
	"github.com/jmoiron/sqlx"
	"log"
)

const (
	limit     = 100
	tableName = "notes"
)

type Repo interface {
	AddNote(ctx context.Context, note *model.Note) (int64, error)
	ListNotes(ctx context.Context) ([]*model.Note, error)
}

type repo struct {
	db sqlx.DB
}

func NewRepo(db sqlx.DB) Repo {
	return &repo{
		db: db,
	}
}

func (r *repo) AddNote(ctx context.Context, note *model.Note) (int64, error) {
	q := sq.Insert(tableName).
		Columns("title", "note_text", "tag").
		Values(note.Title, note.Text, note.Tag).
		RunWith(r.db).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING \"id\"")

	_, _, err := q.ToSql()
	if err != nil {
		return 0, err
	}

	err = q.QueryRowContext(ctx).Scan(&note.Id)
	if err != nil {
		return 0, err
	}

	return note.Id, nil
}

func (r *repo) ListNotes(ctx context.Context) ([]*model.Note, error) {
	var res []*model.Note

	q := sq.Select("*").
		From(tableName).
		RunWith(r.db).
		Limit(limit).
		Offset(1).
		PlaceholderFormat(sq.Dollar)

	rows, err := q.QueryContext(ctx)

	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Fatalf("failed to closing rows: %s", err.Error())
		}
	}(rows)

	for rows.Next() {
		var id, creatingData int64
		var title, text, tag string

		if err = rows.Scan(&id, &title, &text, &tag, &creatingData); err != nil {
			return nil, err
		}

		res = append(res, &model.Note{
			Id:        id,
			Title:     title,
			Text:      text,
			Tag:       tag,
			CreatedAt: creatingData,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}
