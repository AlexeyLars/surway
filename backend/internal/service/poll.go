package service

import (
	"context"
	"fmt"
	"github.com/AlexeyLars/surway-service/internal/config"
	"github.com/AlexeyLars/surway-service/internal/lib/random"
	"github.com/AlexeyLars/surway-service/internal/model"
	"github.com/AlexeyLars/surway-service/internal/storage"
	"log/slog"
	"time"
)

// PollService contains business-logic for poll working
type PollService struct {
	storage storage.Storage
	config  *config.Config
	logger  *slog.Logger
}

func NewPollService(storage storage.Storage, cfg *config.Config, logger *slog.Logger) *PollService {
	return &PollService{
		storage: storage,
		config:  cfg,
		logger:  logger,
	}
}

func (s *PollService) CreatePoll(ctx context.Context, req *model.CreatePollRequest) (*model.CreatePollResponse, error) {
	// Generate short alias as ID
	aliasLen := 7
	pollID := random.NewRandomString(aliasLen)

	now := time.Now()
	expiresAt := now.Add(s.config.Poll.DefaultTTL)

	poll := &model.Poll{
		ID:        pollID,
		Title:     req.Title,
		Options:   req.Options,
		CreatedAt: now,
		ExpiresAt: expiresAt,
	}

	// Save in storage
	if err := s.storage.CreatePoll(ctx, poll, s.config.Poll.DefaultTTL); err != nil {
		s.logger.ErrorContext(ctx, "failed to create poll",
			slog.String("poll_id", pollID),
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to create poll: %w", err)
	}

	s.logger.InfoContext(ctx, "poll created",
		slog.String("poll_id", pollID),
		slog.String("title", req.Title),
		slog.Int("options_count", len(req.Options)),
	)

	// Make response with URL
	baseURL := s.config.Server.BaseURL
	response := &model.CreatePollResponse{
		PollID:     pollID,
		VoteURL:    fmt.Sprintf("%s/api/v1/polls/%s/vote", baseURL, pollID),
		ResultsURL: fmt.Sprintf("%s/api/v1/polls/%s/results", baseURL, pollID),
	}

	return response, nil
}

// Vote register chosen by user option
func (s *PollService) Vote(ctx context.Context, pollID string, req *model.VoteRequest) error {
	if err := s.storage.Vote(ctx, pollID, req.OptionIndices); err != nil {
		if err == storage.ErrPollNotFound {
			s.logger.WarnContext(ctx, "vote for non-existent poll",
				slog.String("poll_id", pollID),
			)
			return err
		}
		if err == storage.ErrInvalidOption {
			s.logger.WarnContext(ctx, "invalid option index",
				slog.String("poll_id", pollID),
				slog.Any("option_indices", req.OptionIndices),
			)
			return err
		}
		if err == storage.ErrDuplicateOption {
			s.logger.WarnContext(ctx, "duplicate option index",
				slog.String("poll_id", pollID),
				slog.Any("option_indices", req.OptionIndices),
			)
			return err
		}

		s.logger.ErrorContext(ctx, "failed to register votes",
			slog.String("poll_id", pollID),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed to register votes: %w", err)
	}

	s.logger.InfoContext(ctx, "votes registered",
		slog.String("poll_id", pollID),
		slog.Any("option_indices", req.OptionIndices),
		slog.Int("votes_count", len(req.OptionIndices)),
	)

	return nil
}

// GetResults get vote results
func (s *PollService) GetResults(ctx context.Context, pollID string) (*model.PollResults, error) {
	results, err := s.storage.GetResults(ctx, pollID)
	if err != nil {
		if err == storage.ErrPollNotFound {
			s.logger.WarnContext(ctx, "results requested for non-existent poll",
				slog.String("poll_id", pollID),
			)
			return nil, err
		}

		s.logger.ErrorContext(ctx, "failed to get results",
			slog.String("poll_id", pollID),
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to get results: %w", err)
	}

	return results, nil
}
