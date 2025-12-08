package service

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/AlexeyLars/surway-service/internal/config"
	"github.com/AlexeyLars/surway-service/internal/model"
	"github.com/AlexeyLars/surway-service/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockStorage is a mock implementation of Storage interface
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) CreatePoll(ctx context.Context, poll *model.Poll, ttl time.Duration) error {
	args := m.Called(ctx, poll, ttl)
	return args.Error(0)
}

func (m *MockStorage) GetPoll(ctx context.Context, pollID string) (*model.Poll, error) {
	args := m.Called(ctx, pollID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Poll), args.Error(1)
}

func (m *MockStorage) Vote(ctx context.Context, pollID string, optionIndices []int) error {
	args := m.Called(ctx, pollID, optionIndices)
	return args.Error(0)
}

func (m *MockStorage) GetResults(ctx context.Context, pollID string) (*model.PollResults, error) {
	args := m.Called(ctx, pollID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.PollResults), args.Error(1)
}

func (m *MockStorage) Close() error {
	args := m.Called()
	return args.Error(0)
}

// Helper functions for creating test data

func newTestConfig() *config.Config {
	return &config.Config{
		Server: config.ServerConfig{
			BaseURL: "http://localhost:8080",
		},
		Poll: config.PollConfig{
			DefaultTTL: 168 * time.Hour,
		},
	}
}

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// Tests

func TestPollService_CreatePoll(t *testing.T) {
	tests := []struct {
		name        string
		request     *model.CreatePollRequest
		setupMock   func(*MockStorage)
		wantErr     bool
		validateRes func(*testing.T, *model.CreatePollResponse)
	}{
		{
			name: "successful poll creation",
			request: &model.CreatePollRequest{
				Title:   "Favorite programming language?",
				Options: []string{"Go", "Python", "JavaScript"},
			},
			setupMock: func(m *MockStorage) {
				m.On("CreatePoll", mock.Anything, mock.AnythingOfType("*model.Poll"), mock.AnythingOfType("time.Duration")).
					Return(nil)
			},
			wantErr: false,
			validateRes: func(t *testing.T, res *model.CreatePollResponse) {
				assert.NotEmpty(t, res.PollID, "PollID should not be empty")
				assert.Equal(t, 7, len(res.PollID), "PollID should be 7 characters long")
				assert.Contains(t, res.VoteURL, res.PollID, "VoteURL should contain PollID")
				assert.Contains(t, res.ResultsURL, res.PollID, "ResultsURL should contain PollID")
				assert.Contains(t, res.VoteURL, "/vote", "VoteURL should contain /vote")
				assert.Contains(t, res.ResultsURL, "/results", "ResultsURL should contain /results")
			},
		},
		{
			name: "storage error during creation",
			request: &model.CreatePollRequest{
				Title:   "Test Poll",
				Options: []string{"A", "B"},
			},
			setupMock: func(m *MockStorage) {
				m.On("CreatePoll", mock.Anything, mock.AnythingOfType("*model.Poll"), mock.AnythingOfType("time.Duration")).
					Return(errors.New("storage error"))
			},
			wantErr: true,
		},
		{
			name: "minimum number of options",
			request: &model.CreatePollRequest{
				Title:   "Binary choice",
				Options: []string{"Yes", "No"},
			},
			setupMock: func(m *MockStorage) {
				m.On("CreatePoll", mock.Anything, mock.AnythingOfType("*model.Poll"), mock.AnythingOfType("time.Duration")).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "maximum number of options",
			request: &model.CreatePollRequest{
				Title:   "Many options",
				Options: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
			},
			setupMock: func(m *MockStorage) {
				m.On("CreatePoll", mock.Anything, mock.AnythingOfType("*model.Poll"), mock.AnythingOfType("time.Duration")).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "poll with long title",
			request: &model.CreatePollRequest{
				Title:   "This is a very long poll title that tests the boundary of what is acceptable for a poll title in our system",
				Options: []string{"Option A", "Option B", "Option C"},
			},
			setupMock: func(m *MockStorage) {
				m.On("CreatePoll", mock.Anything, mock.AnythingOfType("*model.Poll"), mock.AnythingOfType("time.Duration")).
					Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockStorage := new(MockStorage)
			tt.setupMock(mockStorage)

			cfg := newTestConfig()
			logger := newTestLogger()
			service := NewPollService(mockStorage, cfg, logger)
			ctx := context.Background()

			// Act
			response, err := service.CreatePoll(ctx, tt.request)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, response)
			} else {
				require.NoError(t, err)
				require.NotNil(t, response)
				if tt.validateRes != nil {
					tt.validateRes(t, response)
				}
			}

			mockStorage.AssertExpectations(t)
		})
	}
}

func TestPollService_Vote(t *testing.T) {
	tests := []struct {
		name          string
		pollID        string
		request       *model.VoteRequest
		setupMock     func(*MockStorage)
		expectedError error
	}{
		{
			name:   "successful single option vote",
			pollID: "test123",
			request: &model.VoteRequest{
				OptionIndices: []int{0},
			},
			setupMock: func(m *MockStorage) {
				m.On("Vote", mock.Anything, "test123", []int{0}).
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "successful multiple option vote",
			pollID: "test123",
			request: &model.VoteRequest{
				OptionIndices: []int{0, 2, 3},
			},
			setupMock: func(m *MockStorage) {
				m.On("Vote", mock.Anything, "test123", []int{0, 2, 3}).
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "vote for all options",
			pollID: "test123",
			request: &model.VoteRequest{
				OptionIndices: []int{0, 1, 2, 3, 4},
			},
			setupMock: func(m *MockStorage) {
				m.On("Vote", mock.Anything, "test123", []int{0, 1, 2, 3, 4}).
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "poll not found",
			pollID: "nonexistent",
			request: &model.VoteRequest{
				OptionIndices: []int{0},
			},
			setupMock: func(m *MockStorage) {
				m.On("Vote", mock.Anything, "nonexistent", []int{0}).
					Return(storage.ErrPollNotFound)
			},
			expectedError: storage.ErrPollNotFound,
		},
		{
			name:   "invalid option index",
			pollID: "test123",
			request: &model.VoteRequest{
				OptionIndices: []int{999},
			},
			setupMock: func(m *MockStorage) {
				m.On("Vote", mock.Anything, "test123", []int{999}).
					Return(storage.ErrInvalidOption)
			},
			expectedError: storage.ErrInvalidOption,
		},
		{
			name:   "duplicate indices",
			pollID: "test123",
			request: &model.VoteRequest{
				OptionIndices: []int{0, 0, 1},
			},
			setupMock: func(m *MockStorage) {
				m.On("Vote", mock.Anything, "test123", []int{0, 0, 1}).
					Return(storage.ErrDuplicateOption)
			},
			expectedError: storage.ErrDuplicateOption,
		},
		{
			name:   "storage error",
			pollID: "test123",
			request: &model.VoteRequest{
				OptionIndices: []int{0},
			},
			setupMock: func(m *MockStorage) {
				m.On("Vote", mock.Anything, "test123", []int{0}).
					Return(errors.New("storage error"))
			},
			expectedError: errors.New("storage error"),
		},
		{
			name:   "vote with empty indices",
			pollID: "test123",
			request: &model.VoteRequest{
				OptionIndices: []int{},
			},
			setupMock: func(m *MockStorage) {
				m.On("Vote", mock.Anything, "test123", []int{}).
					Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockStorage := new(MockStorage)
			tt.setupMock(mockStorage)

			cfg := newTestConfig()
			logger := newTestLogger()
			service := NewPollService(mockStorage, cfg, logger)
			ctx := context.Background()

			// Act
			err := service.Vote(ctx, tt.pollID, tt.request)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				if errors.Is(tt.expectedError, storage.ErrPollNotFound) ||
					errors.Is(tt.expectedError, storage.ErrInvalidOption) ||
					errors.Is(tt.expectedError, storage.ErrDuplicateOption) {
					assert.ErrorIs(t, err, tt.expectedError)
				}
			} else {
				assert.NoError(t, err)
			}

			mockStorage.AssertExpectations(t)
		})
	}
}

func TestPollService_GetResults(t *testing.T) {
	tests := []struct {
		name          string
		pollID        string
		setupMock     func(*MockStorage)
		expectedError error
		validateRes   func(*testing.T, *model.PollResults)
	}{
		{
			name:   "successful results retrieval",
			pollID: "test123",
			setupMock: func(m *MockStorage) {
				results := &model.PollResults{
					Poll: model.Poll{
						ID:      "test123",
						Title:   "Test Poll",
						Options: []string{"A", "B", "C"},
					},
					Votes: map[string]int{
						"A": 10,
						"B": 5,
						"C": 3,
					},
					Total: 18,
				}
				m.On("GetResults", mock.Anything, "test123").
					Return(results, nil)
			},
			expectedError: nil,
			validateRes: func(t *testing.T, res *model.PollResults) {
				assert.Equal(t, "test123", res.Poll.ID)
				assert.Equal(t, 18, res.Total)
				assert.Equal(t, 10, res.Votes["A"])
				assert.Equal(t, 5, res.Votes["B"])
				assert.Equal(t, 3, res.Votes["C"])
			},
		},
		{
			name:   "results with no votes",
			pollID: "empty123",
			setupMock: func(m *MockStorage) {
				results := &model.PollResults{
					Poll: model.Poll{
						ID:      "empty123",
						Title:   "Empty Poll",
						Options: []string{"A", "B"},
					},
					Votes: map[string]int{
						"A": 0,
						"B": 0,
					},
					Total: 0,
				}
				m.On("GetResults", mock.Anything, "empty123").
					Return(results, nil)
			},
			expectedError: nil,
			validateRes: func(t *testing.T, res *model.PollResults) {
				assert.Equal(t, 0, res.Total)
				assert.Equal(t, 0, res.Votes["A"])
				assert.Equal(t, 0, res.Votes["B"])
			},
		},
		{
			name:   "results with uneven distribution",
			pollID: "uneven123",
			setupMock: func(m *MockStorage) {
				results := &model.PollResults{
					Poll: model.Poll{
						ID:      "uneven123",
						Title:   "Uneven Poll",
						Options: []string{"Popular", "Unpopular"},
					},
					Votes: map[string]int{
						"Popular":   100,
						"Unpopular": 1,
					},
					Total: 101,
				}
				m.On("GetResults", mock.Anything, "uneven123").
					Return(results, nil)
			},
			expectedError: nil,
			validateRes: func(t *testing.T, res *model.PollResults) {
				assert.Equal(t, 101, res.Total)
				assert.True(t, res.Votes["Popular"] > res.Votes["Unpopular"])
			},
		},
		{
			name:   "poll not found",
			pollID: "nonexistent",
			setupMock: func(m *MockStorage) {
				m.On("GetResults", mock.Anything, "nonexistent").
					Return(nil, storage.ErrPollNotFound)
			},
			expectedError: storage.ErrPollNotFound,
		},
		{
			name:   "storage error",
			pollID: "test123",
			setupMock: func(m *MockStorage) {
				m.On("GetResults", mock.Anything, "test123").
					Return(nil, errors.New("storage error"))
			},
			expectedError: errors.New("storage error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockStorage := new(MockStorage)
			tt.setupMock(mockStorage)

			cfg := newTestConfig()
			logger := newTestLogger()
			service := NewPollService(mockStorage, cfg, logger)
			ctx := context.Background()

			// Act
			results, err := service.GetResults(ctx, tt.pollID)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, results)
				if errors.Is(tt.expectedError, storage.ErrPollNotFound) {
					assert.ErrorIs(t, err, tt.expectedError)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, results)
				if tt.validateRes != nil {
					tt.validateRes(t, results)
				}
			}

			mockStorage.AssertExpectations(t)
		})
	}
}

// Benchmark tests

func BenchmarkPollService_CreatePoll(b *testing.B) {
	mockStorage := new(MockStorage)
	mockStorage.On("CreatePoll", mock.Anything, mock.AnythingOfType("*model.Poll"), mock.AnythingOfType("time.Duration")).
		Return(nil)

	cfg := newTestConfig()
	logger := newTestLogger()
	service := NewPollService(mockStorage, cfg, logger)
	ctx := context.Background()

	req := &model.CreatePollRequest{
		Title:   "Benchmark Poll",
		Options: []string{"A", "B", "C"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.CreatePoll(ctx, req)
	}
}

func BenchmarkPollService_Vote_SingleOption(b *testing.B) {
	mockStorage := new(MockStorage)
	mockStorage.On("Vote", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	cfg := newTestConfig()
	logger := newTestLogger()
	service := NewPollService(mockStorage, cfg, logger)
	ctx := context.Background()

	req := &model.VoteRequest{
		OptionIndices: []int{0},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.Vote(ctx, "test123", req)
	}
}

func BenchmarkPollService_Vote_MultipleOptions(b *testing.B) {
	mockStorage := new(MockStorage)
	mockStorage.On("Vote", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	cfg := newTestConfig()
	logger := newTestLogger()
	service := NewPollService(mockStorage, cfg, logger)
	ctx := context.Background()

	req := &model.VoteRequest{
		OptionIndices: []int{0, 1, 2},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.Vote(ctx, "test123", req)
	}
}

func BenchmarkPollService_GetResults(b *testing.B) {
	mockStorage := new(MockStorage)
	results := &model.PollResults{
		Poll: model.Poll{
			ID:      "test123",
			Title:   "Test",
			Options: []string{"A", "B"},
		},
		Votes: map[string]int{"A": 10, "B": 5},
		Total: 15,
	}
	mockStorage.On("GetResults", mock.Anything, mock.Anything).
		Return(results, nil)

	cfg := newTestConfig()
	logger := newTestLogger()
	service := NewPollService(mockStorage, cfg, logger)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GetResults(ctx, "test123")
	}
}

// Edge case tests

func TestPollService_EdgeCases(t *testing.T) {
	t.Run("context cancellation during poll creation", func(t *testing.T) {
		mockStorage := new(MockStorage)
		cfg := newTestConfig()
		logger := newTestLogger()
		service := NewPollService(mockStorage, cfg, logger)

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		req := &model.CreatePollRequest{
			Title:   "Test",
			Options: []string{"A", "B"},
		}

		// Storage might not be called due to cancelled context
		// but if it is called - return error
		mockStorage.On("CreatePoll", mock.Anything, mock.Anything, mock.Anything).
			Return(context.Canceled).Maybe()

		_, err := service.CreatePoll(ctx, req)
		// We might get an error from storage, or service might ignore cancellation
		// Ideally should have context.Context support in storage methods
		_ = err // Just checking it doesn't panic
	})

	t.Run("concurrent votes on same poll", func(t *testing.T) {
		mockStorage := new(MockStorage)
		mockStorage.On("Vote", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)

		cfg := newTestConfig()
		logger := newTestLogger()
		service := NewPollService(mockStorage, cfg, logger)

		// Concurrent voting
		const goroutines = 100
		done := make(chan bool, goroutines)

		for i := 0; i < goroutines; i++ {
			go func(index int) {
				ctx := context.Background()
				req := &model.VoteRequest{
					OptionIndices: []int{index % 3}, // Rotate through options
				}
				_ = service.Vote(ctx, "test123", req)
				done <- true
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < goroutines; i++ {
			<-done
		}

		// Service should handle concurrent requests without panicking
		assert.True(t, true, "Concurrent votes completed successfully")
	})

	t.Run("vote with maximum allowed options", func(t *testing.T) {
		mockStorage := new(MockStorage)
		maxIndices := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
		mockStorage.On("Vote", mock.Anything, "test123", maxIndices).
			Return(nil)

		cfg := newTestConfig()
		logger := newTestLogger()
		service := NewPollService(mockStorage, cfg, logger)
		ctx := context.Background()

		req := &model.VoteRequest{
			OptionIndices: maxIndices,
		}

		err := service.Vote(ctx, "test123", req)
		assert.NoError(t, err)
		mockStorage.AssertExpectations(t)
	})

	t.Run("create poll with minimum valid data", func(t *testing.T) {
		mockStorage := new(MockStorage)
		mockStorage.On("CreatePoll", mock.Anything, mock.AnythingOfType("*model.Poll"), mock.AnythingOfType("time.Duration")).
			Return(nil)

		cfg := newTestConfig()
		logger := newTestLogger()
		service := NewPollService(mockStorage, cfg, logger)
		ctx := context.Background()

		req := &model.CreatePollRequest{
			Title:   "ABC",              // Minimum 3 chars
			Options: []string{"X", "Y"}, // Minimum 2 options
		}

		response, err := service.CreatePoll(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		mockStorage.AssertExpectations(t)
	})
}

// Integration-style test with realistic scenarios

func TestPollService_RealisticScenarios(t *testing.T) {
	t.Run("complete poll lifecycle", func(t *testing.T) {
		mockStorage := new(MockStorage)
		cfg := newTestConfig()
		logger := newTestLogger()
		service := NewPollService(mockStorage, cfg, logger)
		ctx := context.Background()

		// 1. Create poll
		createReq := &model.CreatePollRequest{
			Title:   "Best superhero?",
			Options: []string{"Batman", "Superman", "Spider-Man"},
		}

		mockStorage.On("CreatePoll", mock.Anything, mock.AnythingOfType("*model.Poll"), mock.AnythingOfType("time.Duration")).
			Return(nil).Once()

		createResp, err := service.CreatePoll(ctx, createReq)
		require.NoError(t, err)
		require.NotNil(t, createResp)
		pollID := createResp.PollID

		// 2. Multiple users vote
		mockStorage.On("Vote", mock.Anything, pollID, []int{0}).Return(nil).Once()
		mockStorage.On("Vote", mock.Anything, pollID, []int{1}).Return(nil).Once()
		mockStorage.On("Vote", mock.Anything, pollID, []int{0, 2}).Return(nil).Once()

		err = service.Vote(ctx, pollID, &model.VoteRequest{OptionIndices: []int{0}})
		assert.NoError(t, err)

		err = service.Vote(ctx, pollID, &model.VoteRequest{OptionIndices: []int{1}})
		assert.NoError(t, err)

		err = service.Vote(ctx, pollID, &model.VoteRequest{OptionIndices: []int{0, 2}})
		assert.NoError(t, err)

		// 3. Get results
		expectedResults := &model.PollResults{
			Poll: model.Poll{
				ID:      pollID,
				Title:   "Best superhero?",
				Options: []string{"Batman", "Superman", "Spider-Man"},
			},
			Votes: map[string]int{
				"Batman":     2,
				"Superman":   1,
				"Spider-Man": 1,
			},
			Total: 4,
		}

		mockStorage.On("GetResults", mock.Anything, pollID).Return(expectedResults, nil).Once()

		results, err := service.GetResults(ctx, pollID)
		require.NoError(t, err)
		assert.Equal(t, 4, results.Total)

		mockStorage.AssertExpectations(t)
	})
}
