package services

import (
	"context"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
	"github.com/champion19/flighthours-api/platform/logger"
)

// MessageService handles message-related business logic
type MessageService struct {
	repository output.MessageRepository
	logger     logger.Logger
}

// NewMessageService creates a new message service
func NewMessageService(repository output.MessageRepository, log logger.Logger) *MessageService {
	return &MessageService{
		repository: repository,
		logger:     log,
	}
}

// BeginTx starts a new transaction
func (s *MessageService) BeginTx(ctx context.Context) (output.Tx, error) {
	return s.repository.BeginTx(ctx)
}

// ValidateMessage validates a message before saving
func (s *MessageService) ValidateMessage(ctx context.Context, message domain.Message) error {
	s.logger.Debug(logger.LogMessageValidation, "code", message.Code)

	if err := message.Validate(); err != nil {
		s.logger.Error(logger.LogMessageValidationError, "code", message.Code, "error", err)
		return err
	}

	// Check if code already exists (for create operations)
	if message.ID == "" {
		existing, err := s.repository.GetByCode(ctx, message.Code)
		if err == nil && existing != nil {
			s.logger.Warn(logger.LogMessageCodeDuplicate, "code", message.Code)
			return domain.ErrMessageCodeDuplicate
		}
	}

	s.logger.Debug(logger.LogMessageValidationOK, "code", message.Code)
	return nil
}

// GetMessageByID retrieves a message by ID (includes inactive messages for CRUD operations)
func (s *MessageService) GetMessageByID(ctx context.Context, id string) (*domain.Message, error) {
	s.logger.Debug(logger.LogMessageGet, "id", id)

	message, err := s.repository.GetByID(ctx, id)
	if err != nil {
		s.logger.Error(logger.LogMessageGetError, "id", id, "error", err)
		return nil, err
	}

	if message == nil {
		s.logger.Warn(logger.LogMessageGetError, "id", id, "error", "not found")
		return nil, domain.ErrMessageNotFound
	}

	s.logger.Debug(logger.LogMessageGetOK, message.ToLogger())
	return message, nil
}

// GetMessageByCode retrieves a message by code
func (s *MessageService) GetMessageByCode(ctx context.Context, code string) (*domain.Message, error) {
	s.logger.Debug(logger.LogMessageGet, "code", code)

	message, err := s.repository.GetByCode(ctx, code)
	if err != nil {
		s.logger.Error(logger.LogMessageGetError, "code", code, "error", err)
		return nil, err
	}

	s.logger.Debug(logger.LogMessageGetOK, message.ToLogger())
	return message, nil
}

// ListMessages lists messages with optional filters
func (s *MessageService) ListMessages(ctx context.Context, filters map[string]interface{}) ([]domain.Message, error) {
	s.logger.Debug(logger.LogMessageList, "filters", filters)

	// ListMessages with filters - use GetByType or GetByModule based on filters
	messages, err := s.repository.GetAllActive(ctx)
	if err != nil {
		s.logger.Error(logger.LogMessageListError, "error", err)
		return nil, err
	}

	s.logger.Debug(logger.LogMessageListOK, "count", len(messages))
	return messages, nil
}

// ListActiveMessages lists only active messages
func (s *MessageService) ListActiveMessages(ctx context.Context) ([]domain.Message, error) {
	s.logger.Debug(logger.LogMessageList, "filter", "active_only")

	messages, err := s.repository.GetAllActive(ctx)
	if err != nil {
		s.logger.Error(logger.LogMessageListError, "error", err)
		return nil, err
	}

	s.logger.Debug(logger.LogMessageListOK, "count", len(messages))
	return messages, nil
}

// SaveMessageToDB saves a message to the database
func (s *MessageService) SaveMessageToDB(ctx context.Context, tx output.Tx, message domain.Message) error {
	s.logger.Info(logger.LogMessageCreate, message.ToLogger())

	err := s.repository.SaveMessage(ctx, tx, message)
	if err != nil {
		s.logger.Error(logger.LogMessageCreateError, message.ToLogger(), "error", err)
		return err
	}

	s.logger.Success(logger.LogMessageCreateOK, message.ToLogger())
	return nil
}

// UpdateMessageInDB updates a message in the database
func (s *MessageService) UpdateMessageInDB(ctx context.Context, tx output.Tx, message domain.Message) error {
	s.logger.Info(logger.LogMessageUpdate, message.ToLogger())

	existingMsg, err := s.repository.GetByID(ctx, message.ID)
	if err != nil {
		return domain.ErrMessageCannotUpdate
	}

	if existingMsg == nil || existingMsg.Code != message.Code {
		return domain.ErrMessageNotFound
	}

	err = s.repository.UpdateMessage(ctx, tx, message)
	if err != nil {
		s.logger.Error(logger.LogMessageUpdateError, message.ToLogger(), "error", err)
		return err
	}

	s.logger.Success(logger.LogMessageUpdateOK, message.ToLogger())
	return nil
}

// DeleteMessageFromDB deletes a message from the database
func (s *MessageService) DeleteMessageFromDB(ctx context.Context, tx output.Tx, id string) error {
	s.logger.Info(logger.LogMessageDelete, "id", id)

	err := s.repository.DeleteMessage(ctx, tx, id)
	if err != nil {
		s.logger.Error(logger.LogMessageDeleteError, "id", id, "error", err)
		return err
	}

	s.logger.Success(logger.LogMessageDeleteOK, "id", id)
	return nil
}
