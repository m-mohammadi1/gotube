package channelrepo

import (
	"context"
	"database/sql"
	"gotube/pkg/model"
)

type PostgresRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) UpdateOrCreate(ctx context.Context, channel model.Channel) (*model.Channel, error) {
	query := `
		INSERT INTO channels (user_id, youtube_id, title, added_at, token)
        VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (youtube_id) DO UPDATE 
		SET user_id = EXCLUDED.user_id, title = EXCLUDED.title, added_at = EXCLUDED.added_at, token = EXCLUDED.token
		RETURNING id
        `

	err := r.db.QueryRow(query, channel.UserID, channel.YoutubeID, channel.Title, channel.AddedAt, channel.Token).Scan(&channel.ID)

	if err != nil {
		return nil, err
	}

	return &channel, nil
}
