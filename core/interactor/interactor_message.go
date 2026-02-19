package interactor

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/helpers"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/logger"
)

type MessageInteractor struct {
	service input.MessageService
}

func NewMessageInteractor(service input.MessageService) *MessageInteractor {
	return &MessageInteractor{
		service: service,
	}
}

func (i *MessageInteractor) CreateMessage(ctx context.Context, message domain.Message) (result *domain.Message, err error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogMessageCreate, message.ToLogger())

	if err = i.service.ValidateMessage(ctx, message); err != nil {
		log.Error(logger.LogMessageInteractorCreateStep1Error, "error", err)
		return
	}
	log.Success(logger.LogMessageInteractorCreateStep1OK)

	err = helpers.RunWithTx(ctx, i.service, log, logger.LogMessageInteractorCreateStep2Error,
		func(ctx context.Context, tx output.Tx) error {
			return i.service.SaveMessageToDB(ctx, tx, message)
		})
	if err != nil {
		return
	}

	result = &message
	log.Success(logger.LogMessageInteractorCreateComplete, message.ToLogger())

	err = nil
	return
}

func (i *MessageInteractor) UpdateMessage(ctx context.Context, message domain.Message) (result *domain.Message, err error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogMessageUpdate, message.ToLogger())

	_, err = i.service.GetMessageByID(ctx, message.ID)
	if err != nil {
		log.Error(logger.LogMessageInteractorUpdateStep1Error, "error", err)
		return nil, err
	}
	log.Success(logger.LogMessageInteractorUpdateStep1OK)

	if err = i.service.ValidateMessage(ctx, message); err != nil {
		log.Error(logger.LogMessageInteractorUpdateStep2Error, "error", err)
		return
	}
	log.Success(logger.LogMessageInteractorUpdateStep2OK)

	err = helpers.RunWithTx(ctx, i.service, log, logger.LogMessageInteractorUpdateStep3Error,
		func(ctx context.Context, tx output.Tx) error {
			return i.service.UpdateMessageInDB(ctx, tx, message)
		})
	if err != nil {
		return
	}

	result = &message
	log.Success(logger.LogMessageInteractorUpdateComplete, message.ToLogger())

	log.Info(logger.LogMessageCacheRefresh)

	err = nil
	return
}

func (i *MessageInteractor) DeleteMessage(ctx context.Context, id string) (err error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogMessageDelete, "id", id)

	_, err = i.service.GetMessageByID(ctx, id)
	if err != nil {
		log.Error(logger.LogMessageInteractorDeleteStep1Error, "error", err)
		return err
	}
	log.Success(logger.LogMessageInteractorDeleteStep1OK)

	err = helpers.RunWithTx(ctx, i.service, log, logger.LogMessageInteractorDeleteStep2Error,
		func(ctx context.Context, tx output.Tx) error {
			return i.service.DeleteMessageFromDB(ctx, tx, id)
		})
	if err != nil {
		return err
	}

	log.Success(logger.LogMessageInteractorDeleteComplete, "id", id)

	log.Info(logger.LogMessageCacheRefresh)

	err = nil
	return
}

func (i *MessageInteractor) GetMessageByID(ctx context.Context, id string) (*domain.Message, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Debug(logger.LogMessageGet, "id", id)

	message, err := i.service.GetMessageByID(ctx, id)
	if err != nil {
		log.Error(logger.LogMessageGetError, "id", id, "error", err)
		return nil, err
	}

	log.Debug(logger.LogMessageGetOK, message.ToLogger())
	return message, nil
}

func (i *MessageInteractor) GetMessageByCode(ctx context.Context, code string) (*domain.Message, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Debug(logger.LogMessageGet, "code", code)

	message, err := i.service.GetMessageByCode(ctx, code)
	if err != nil {
		log.Error(logger.LogMessageGetError, "code", code, "error", err)
		return nil, err
	}

	log.Debug(logger.LogMessageGetOK, message.ToLogger())
	return message, nil
}

func (i *MessageInteractor) ListMessages(ctx context.Context, filters map[string]interface{}) ([]domain.Message, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Debug(logger.LogMessageList, "filters", filters)

	messages, err := i.service.ListMessages(ctx, filters)
	if err != nil {
		log.Error(logger.LogMessageListError, "error", err)
		return nil, err
	}

	log.Debug(logger.LogMessageListOK, "count", len(messages))
	return messages, nil
}

func (i *MessageInteractor) ListActiveMessages(ctx context.Context) ([]domain.Message, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Debug(logger.LogMessageList, "filter", "active_only")

	messages, err := i.service.ListActiveMessages(ctx)
	if err != nil {
		log.Error(logger.LogMessageListError, "error", err)
		return nil, err
	}

	log.Debug(logger.LogMessageListOK, "count", len(messages))
	return messages, nil
}
