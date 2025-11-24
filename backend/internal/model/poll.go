package model

import (
	"time"
)

type Poll struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Options   []string  `json:"options"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type CreatePollRequest struct {
	Title   string   `json:"title" binding:"required,min=3,max=200"`
	Options []string `json:"options" binding:"required,min=2,max=10,dive,required,min=1,max=100"`
}

type CreatePollResponse struct {
	PollID     string `json:"poll_id"`
	VoteURL    string `json:"vote_url"`
	ResultsURL string `json:"results_url"`
}

type VoteRequest struct {
	OptionIndices []int `json:"option_indices" binding:"required,min=1,dive,min=0"`
}

type VoteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type PollResults struct {
	Poll  Poll           `json:"poll"`
	Votes map[string]int `json:"votes"` // option ->  count ...
	Total int            `json:"total"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
