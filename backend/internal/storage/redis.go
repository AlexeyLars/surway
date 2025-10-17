package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AlexeyLars/surway-service/internal/model"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	// ErrPollNotFound returns when poll not found
	ErrPollNotFound = errors.New("poll not found")

	// ErrInvalidOption returns when trying vote for non existed option
	ErrInvalidOption = errors.New("invalid option index")

	// ErrDuplicateOption returns when options indices not unique
	ErrDuplicateOption = errors.New("options indices not unique")
)

// Storage defines interface for working with polls storage
type Storage interface {
	CreatePoll(ctx context.Context, poll *model.Poll, ttl time.Duration) error
	GetPoll(ctx context.Context, pollID string) (*model.Poll, error)
	Vote(ctx context.Context, pollID string, optionIndices []int) error
	GetResults(ctx context.Context, pollID string) (*model.PollResults, error)
	Close() error
}

// RedisStorage realises Storage for Redis
type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(client *redis.Client) *RedisStorage {
	return &RedisStorage{
		client: client,
	}
}

// Keys for Redis
func pollInfoKey(pollID string) string {
	return fmt.Sprintf("poll:%s:info", pollID)
}

func pollVotesKey(pollID string) string {
	return fmt.Sprintf("poll:%s:votes", pollID)
}

// CreatePoll saves new poll in Redis
func (s *RedisStorage) CreatePoll(ctx context.Context, poll *model.Poll, ttl time.Duration) error {
	pollData, err := json.Marshal(poll)
	if err != nil {
		return fmt.Errorf("failed to marshal poll: %w", err)
	}

	pipe := s.client.Pipeline()

	// Save poll info
	pipe.Set(ctx, pollInfoKey(poll.ID), pollData, ttl)

	// Initialize votes counters as zeros
	votesData := make(map[string]interface{})
	for i := range poll.Options {
		votesData[fmt.Sprintf("%d", i)] = 0
	}
	pipe.HSet(ctx, pollVotesKey(poll.ID), votesData)
	pipe.Expire(ctx, pollVotesKey(poll.ID), ttl)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create poll: %w", err)
	}

	return nil
}

// GetPoll gets poll info
func (s *RedisStorage) GetPoll(ctx context.Context, pollID string) (*model.Poll, error) {
	data, err := s.client.Get(ctx, pollInfoKey(pollID)).Result()
	if err == redis.Nil {
		return nil, ErrPollNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get poll: %w", err)
	}

	var poll model.Poll
	if err := json.Unmarshal([]byte(data), &poll); err != nil {
		return nil, fmt.Errorf("failed to unmarshal poll: %w", err)
	}

	return &poll, nil
}

func (s *RedisStorage) Vote(ctx context.Context, pollID string, optionIndices []int) error {
	// Check if poll exists
	exists, err := s.client.Exists(ctx, pollInfoKey(pollID)).Result()
	if err != nil {
		return fmt.Errorf("failed to check poll existence: %w", err)
	}
	if exists == 0 {
		return ErrPollNotFound
	}

	// Get poll, validate indices and check for duplicates
	poll, err := s.GetPoll(ctx, pollID)
	if err != nil {
		return err
	}

	for _, idx := range optionIndices {
		if idx < 0 || idx >= len(poll.Options) {
			return ErrInvalidOption
		}
	}

	seen := make(map[int]bool)
	for _, idx := range optionIndices {
		if seen[idx] {
			return ErrDuplicateOption
		}
		seen[idx] = true
	}

	// Atomic counters increase
	pipe := s.client.Pipeline()
	for _, idx := range optionIndices {
		field := fmt.Sprintf("%d", idx)
		pipe.HIncrBy(ctx, pollVotesKey(pollID), field, 1)
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to register votes: %w", err)
	}

	return nil
}

func (s *RedisStorage) GetResults(ctx context.Context, pollID string) (*model.PollResults, error) {
	poll, err := s.GetPoll(ctx, pollID)
	if err != nil {
		return nil, err
	}

	votesMap, err := s.client.HGetAll(ctx, pollVotesKey(pollID)).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get votes: %w", err)
	}

	// Convert votes into map
	votes := make(map[string]int)
	total := 0

	for i, option := range poll.Options {
		field := fmt.Sprintf("%d", i)
		count := 0
		if countStr, ok := votesMap[field]; ok {
			fmt.Sscanf(countStr, "%d", &count)
		}
		votes[option] = count
		total += count
	}

	return &model.PollResults{
		Poll:  *poll,
		Votes: votes,
		Total: total,
	}, nil
}

func (s *RedisStorage) Close() error {
	return s.client.Close()
}
