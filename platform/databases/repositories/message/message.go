package message

import (
	"time"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

type Message struct {
	ID string `json:"id"`
	Code string `json:"code"`
	Type string `json:"type"`
	Category string `json:"category"`
	Module string `json:"module"`
	Title string `json:"title"`
	Content string `json:"content"`
	Active bool `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (r *Message) ToDomain() domain.Message {
	return domain.Message{
		ID:        r.ID,
		Code:      r.Code,
		Type:      domain.MessageType(r.Type),
		Category:  r.Category,
		Module:    r.Module,
		Title:     r.Title,
		Content:   r.Content,
		Active:    r.Active,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func FromDomain(domainMessage domain.Message) Message {
	return Message{
		ID:        domainMessage.ID,
		Code:      domainMessage.Code,
		Type:      string(domainMessage.Type),
		Category:  domainMessage.Category,
		Module:    domainMessage.Module,
		Title:     domainMessage.Title,
		Content:   domainMessage.Content,
		Active:    domainMessage.Active,
		CreatedAt: domainMessage.CreatedAt,
		UpdatedAt: domainMessage.UpdatedAt,
	}
}
