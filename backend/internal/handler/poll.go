package handler

import (
	"errors"
	"github.com/AlexeyLars/surway-service/internal/model"
	"github.com/AlexeyLars/surway-service/internal/service"
	"github.com/AlexeyLars/surway-service/internal/storage"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

// PollHandler processes HTTP votes requests
type PollHandler struct {
	service *service.PollService
	logger  *slog.Logger
}

// NewPollHandler creates new PollHandler
func NewPollHandler(service *service.PollService, logger *slog.Logger) *PollHandler {
	return &PollHandler{
		service: service,
		logger:  logger,
	}
}

// CreatePoll godoc
// @Summary      Create poll
// @Description  Create new poll with given name and options
// @Tags         polls
// @Accept       json
// @Produce      json
// @Param        request body model.CreatePollRequest true "Poll data"
// @Success      201 {object} model.CreatePollResponse
// @Failure      400 {object} model.ErrorResponse
// @Failure      500 {object} model.ErrorResponse
// @Router       /polls [post]
func (h *PollHandler) CreatePoll(c *gin.Context) {
	var req model.CreatePollRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WarnContext(c.Request.Context(), "invalid request body",
			slog.String("error", err.Error()),
		)
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	response, err := h.service.CreatePoll(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to create poll",
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// Vote godoc
// @Summary      Vote
// @Description  Register vote for given option
// @Tags         polls
// @Accept       json
// @Produce      json
// @Param        id path string true "Poll ID"
// @Param        request body model.VoteRequest true "Option index"
// @Success      200 {object} model.VoteResponse
// @Failure      400 {object} model.ErrorResponse
// @Failure      404 {object} model.ErrorResponse
// @Failure      500 {object} model.ErrorResponse
// @Router       /polls/{id}/vote [post]
func (h *PollHandler) Vote(c *gin.Context) {
	pollID := c.Param("id")

	var req model.VoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WarnContext(c.Request.Context(), "invalid vote request",
			slog.String("poll_id", pollID),
			slog.String("error", err.Error()),
		)
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	err := h.service.Vote(c.Request.Context(), pollID, &req)
	if err != nil {
		if errors.Is(err, storage.ErrPollNotFound) {
			c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error:   "poll_not_found",
				Message: "Poll not found or expired",
			})
			return
		}
		if errors.Is(err, storage.ErrInvalidOption) {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Error:   "invalid_option",
				Message: "Invalid option index",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to register vote",
		})
		return
	}

	c.JSON(http.StatusOK, model.VoteResponse{
		Success: true,
		Message: "Vote registered successfully",
	})
}

// GetResults godoc
// @Summary      Get results
// @Description  Return poll results with option's vote counts
// @Tags         polls
// @Produce      json
// @Param        id path string true "Poll ID"
// @Success      200 {object} model.PollResults
// @Failure      404 {object} model.ErrorResponse
// @Failure      500 {object} model.ErrorResponse
// @Router       /polls/{id}/results [get]
func (h *PollHandler) GetResults(c *gin.Context) {
	pollID := c.Param("id")

	results, err := h.service.GetResults(c.Request.Context(), pollID)
	if err != nil {
		if errors.Is(err, storage.ErrPollNotFound) {
			c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error:   "poll_not_found",
				Message: "Poll not found or expired",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to get results",
		})
		return
	}

	c.JSON(http.StatusOK, results)
}

// HealthCheck godoc
// @Summary      Health check
// @Description  Check service heath state
// @Tags         health
// @Produce      json
// @Success      200 {object} map[string]string
// @Router       /health [get]
func (h *PollHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
