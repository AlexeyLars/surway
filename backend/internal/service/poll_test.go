package service

import (
	"context"
	"errors"
	"github.com/AlexeyLars/surway-service/internal/config"
	"github.com/AlexeyLars/surway-service/internal/model"
	"github.com/AlexeyLars/surway-service/internal/storage"
	"log/slog"
	"os"
	"testing"
	"time"
)

// MockStorage - mock for testing
type MockStorage struct {
	CreatePollFunc func(ctx context.Context, poll *model.Poll, ttl time.Duration) error
	GetPollFunc    func(ctx context.Context, pollID string) (*model.Poll, error)
	VoteFunc       func(ctx context.Context, pollID string, optionIndex int) error
	GetResultsFunc func(ctx context.Context, pollID string) (*model.PollResults, error)
}

func (m *MockStorage) CreatePoll(ctx context.Context, poll *model.Poll, ttl time.Duration) error {
	if m.CreatePollFunc != nil {
		return m.CreatePollFunc(ctx, poll, ttl)
	}
	return nil
}

func (m *MockStorage) GetPoll(ctx context.Context, pollID string) (*model.Poll, error) {
	if m.GetPollFunc != nil {
		return m.GetPollFunc(ctx, pollID)
	}
	return nil, storage.ErrPollNotFound
}

func (m *MockStorage) Vote(ctx context.Context, pollID string, optionIndex int) error {
	if m.VoteFunc != nil {
		return m.VoteFunc(ctx, pollID, optionIndex)
	}
	return nil
}

func (m *MockStorage) GetResults(ctx context.Context, pollID string) (*model.PollResults, error) {
	if m.GetResultsFunc != nil {
		return m.GetResultsFunc(ctx, pollID)
	}
	return nil, storage.ErrPollNotFound
}

func (m *MockStorage) Close() error {
	return nil
}

func TestPollService_CreatePoll(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError, // only errors in tests
	}))

	cfg := &config.Config{
		Server: config.ServerConfig{
			BaseURL: "http://localhost:8080",
		},
		Poll: config.PollConfig{
			DefaultTTL: 168 * time.Hour,
		},
	}

	tests := []struct {
		name      string
		request   *model.CreatePollRequest
		mockError error
		wantError bool
	}{
		{
			name: "Successful poll creation",
			request: &model.CreatePollRequest{
				Title:   "Test Poll",
				Options: []string{"Option 1", "Option 2"},
			},
			mockError: nil,
			wantError: false,
		},
		{
			name: "Error with saving",
			request: &model.CreatePollRequest{
				Title:   "Test Poll",
				Options: []string{"Option 1", "Option 2"},
			},
			mockError: errors.New("storage error"),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockStorage{
				CreatePollFunc: func(ctx context.Context, poll *model.Poll, ttl time.Duration) error {
					return tt.mockError
				},
			}

			service := NewPollService(mockStorage, cfg, logger)
			ctx := context.Background()

			resp, err := service.CreatePoll(ctx, tt.request)

			if tt.wantError {
				if err == nil {
					t.Error("Expected error, actual success")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if resp.PollID == "" {
				t.Error("poll_id can not be empty")
			}

			if resp.VoteURL == "" {
				t.Error("vote_url can not be empty")
			}

			if resp.ResultsURL == "" {
				t.Error("results_url can not be empty")
			}
		})
	}
}

func TestPollService_Vote(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	cfg := &config.Config{
		Poll: config.PollConfig{
			DefaultTTL: 168 * time.Hour,
		},
	}

	tests := []struct {
		name      string
		pollID    string
		request   *model.VoteRequest
		mockError error
		wantError error
	}{
		{
			name:   "Succesful vote",
			pollID: "test-poll-id",
			request: &model.VoteRequest{
				OptionIndex: 0,
			},
			mockError: nil,
			wantError: nil,
		},
		{
			name:   "Poll not found",
			pollID: "non-existent",
			request: &model.VoteRequest{
				OptionIndex: 0,
			},
			mockError: storage.ErrPollNotFound,
			wantError: storage.ErrPollNotFound,
		},
		{
			name:   "Non-valid poll option",
			pollID: "test-poll-id",
			request: &model.VoteRequest{
				OptionIndex: 999,
			},
			mockError: storage.ErrInvalidOption,
			wantError: storage.ErrInvalidOption,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockStorage{
				VoteFunc: func(ctx context.Context, pollID string, optionIndex int) error {
					return tt.mockError
				},
			}

			service := NewPollService(mockStorage, cfg, logger)
			ctx := context.Background()

			err := service.Vote(ctx, tt.pollID, tt.request)

			if tt.wantError != nil {
				if !errors.Is(err, tt.wantError) {
					t.Errorf("Expected error %v, actual %v", tt.wantError, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
