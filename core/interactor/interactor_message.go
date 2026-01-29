package interactor

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
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

	tx, err := i.service.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogMessageInteractorCreateStep2Error, "error", err)
		return
	}
	log.Success(logger.LogMessageInteractorCreateStep2OK)

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogMessageInteractorRollbackError, "error", rbErr)
			} else {
				log.Warn(logger.LogMessageInteractorRollbackOK, "error", err)
			}
		}
	}()

	if err = i.service.SaveMessageToDB(ctx, tx, message); err != nil {
		log.Error(logger.LogMessageInteractorCreateStep3Error, "error", err)
		return
	}
	log.Success(logger.LogMessageInteractorCreateStep3OK)

	if err = tx.Commit(); err != nil {
		log.Error(logger.LogMessageInteractorCreateCommitErr, "error", err)
		return
	}
	log.Success(logger.LogMessageInteractorCreateCommitOK)

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

	tx, err := i.service.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogMessageInteractorUpdateStep3Error, "error", err)
		return
	}
	log.Success(logger.LogMessageInteractorUpdateStep3OK)

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogMessageInteractorRollbackError,
					"rollback error", rbErr, "original error", err)
			} else {
				log.Warn(logger.LogMessageInteractorRollbackOK)
			}
		}
	}()

	if err = i.service.UpdateMessageInDB(ctx, tx, message); err != nil {
		log.Error(logger.LogMessageInteractorUpdateStep4Error, "error", err)
		return
	}
	log.Success(logger.LogMessageInteractorUpdateStep4OK)

	if err = tx.Commit(); err != nil {
		log.Error(logger.LogMessageInteractorUpdateCommitErr, "error", err)
		return
	}
	log.Success(logger.LogMessageInteractorUpdateCommitOK)

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

	tx, err := i.service.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogMessageInteractorDeleteStep2Error, "error", err)
		return err
	}
	log.Success(logger.LogMessageInteractorDeleteStep2OK)

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogMessageInteractorRollbackError,
					"rollback_error", rbErr,
					"original_error", err)
			} else {
				log.Warn(logger.LogMessageInteractorRollbackOK,
					"original_error", err)
			}
		}
	}()

	if err = i.service.DeleteMessageFromDB(ctx, tx, id); err != nil {
		log.Error(logger.LogMessageInteractorDeleteStep3Error, "error", err)
		return err
	}
	log.Success(logger.LogMessageInteractorDeleteStep3OK)

	if err = tx.Commit(); err != nil {
		log.Error(logger.LogMessageInteractorDeleteCommitErr, "error", err)
		return err
	}
	log.Success(logger.LogMessageInteractorDeleteCommitOK)

	log.Success("Mensaje eliminado exitosamente", "id", id)

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
