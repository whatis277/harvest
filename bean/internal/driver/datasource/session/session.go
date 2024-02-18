package session

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/whatis277/harvest/bean/internal/entity/model"

	"github.com/whatis277/harvest/bean/internal/usecase/interfaces"

	"github.com/whatis277/harvest/bean/internal/driver/redis"
)

type dataSource struct {
	cache *redis.Cache
}

func New(cache *redis.Cache) interfaces.SessionDataSource {
	return &dataSource{
		cache: cache,
	}
}

func (ds *dataSource) Create(
	userID string,
	hashedToken string,
	duration time.Duration,
) (*model.Session, error) {
	return ds.doCreate(userID, hashedToken, duration, 0)
}

func (ds *dataSource) doCreate(
	userID string,
	hashedToken string,
	duration time.Duration,
	attempts int,
) (*model.Session, error) {
	if attempts > 3 {
		return nil, fmt.Errorf("failed to generate session: max attempts exceeded")
	}

	rand, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to generate random session id: %w", err)
	}

	id := rand.String()

	existing, err := ds.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find existing session: %w", err)
	}

	if existing != nil {
		return ds.doCreate(userID, hashedToken, duration, attempts+1)
	}

	session := &model.Session{
		ID:          id,
		UserID:      userID,
		HashedToken: hashedToken,
	}

	err = ds.cache.Client.Set(context.Background(), id, session, duration).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to set session in cache: %w", err)
	}

	return session, nil
}

func (ds *dataSource) FindByID(id string) (*model.Session, error) {
	session := &model.Session{}

	err := ds.cache.Client.Get(context.Background(), id).Scan(session)
	if err == redis.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get session from cache: %w", err)
	}

	return session, nil
}

func (ds *dataSource) Refresh(id string, duration time.Duration) error {
	err := ds.cache.Client.Expire(context.Background(), id, duration).Err()
	if err != nil {
		return fmt.Errorf("failed to refresh session in cache: %w", err)
	}

	return nil
}

func (ds *dataSource) Delete(id string) error {
	err := ds.cache.Client.Del(context.Background(), id).Err()
	if err != nil {
		return fmt.Errorf("failed to delete session from cache: %w", err)
	}

	return nil
}
