package token

import (
	"context"
	"fmt"

	"github.com/whatis277/harvest/bean/internal/entity/model"

	"github.com/whatis277/harvest/bean/internal/usecase/interfaces"

	"github.com/whatis277/harvest/bean/internal/driver/postgres"
)

type dataSource struct {
	db *postgres.DB
}

func New(db *postgres.DB) interfaces.LoginTokenDataSource {
	return &dataSource{
		db: db,
	}
}

func (ds *dataSource) Create(email string, hashedToken string) (*model.LoginToken, error) {
	token := &model.LoginToken{}

	err := ds.db.Pool.
		QueryRow(
			context.Background(),
			("INSERT INTO login_tokens (email, hashed_token)"+
				" VALUES ($1, $2)"+
				" ON CONFLICT (email) DO UPDATE"+
				" SET hashed_token = $2"+
				" RETURNING *"),
			email, hashedToken,
		).
		Scan(&token.ID, &token.Email, &token.HashedToken, &token.CreatedAt, &token.ExpiresAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}

	return token, nil
}

func (ds *dataSource) FindUnexpired(id string) (*model.LoginToken, error) {
	token := &model.LoginToken{}

	err := ds.db.Pool.
		QueryRow(
			context.Background(),
			("SELECT * FROM login_tokens"+
				" WHERE id = $1 AND expires_at > NOW()"),
			id,
		).
		Scan(&token.ID, &token.Email, &token.HashedToken, &token.CreatedAt, &token.ExpiresAt)

	if err == postgres.ErrNowRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find unexpired token: %w", err)
	}

	return token, nil
}

func (ds *dataSource) Delete(id string) error {
	_, err := ds.db.Pool.
		Exec(
			context.Background(),
			"DELETE FROM login_tokens WHERE id = $1",
			id,
		)

	if err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}

	return nil
}
