package handlers
import (

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)
type MessageRequest struct {
	Code     string             `json:"code" binding:"omitempty"`
	Type     domain.MessageType `json:"type" binding:"omitempty"`
	Category string             `json:"category" binding:"omitempty"`
	Module   string             `json:"module" binding:"omitempty"`
	Title    string             `json:"title" binding:"omitempty"`
	Content  string             `json:"content" binding:"omitempty"`
	Active   bool               `json:"active"`
}


type MessageResponse struct {
	ID       string             `json:"id"`
	Code     string             `json:"code"`
	Type     domain.MessageType `json:"type"`
	Category string             `json:"category"`
	Module   string             `json:"module"`
	Title    string             `json:"title"`
	Content  string             `json:"content"`
	Active   bool               `json:"active"`
	Links    []Link             `json:"_links,omitempty"`
}

type MessageListResponse struct {
	Messages []MessageResponse `json:"messages"`
	Count    int               `json:"count"`
	Links    []Link            `json:"_links,omitempty"`
}

type MessageCreatedResponse struct {
	ID    string `json:"id"`
	Links []Link `json:"_links"`

}

type MessageUpdatedResponse struct {
	Links []Link `json:"_links"`

}

type MessageDeletedResponse struct {
}


type CacheReloadResponse struct {
	Success     bool   `json:"success"`
	BeforeCount int    `json:"before_count"`
	AfterCount  int    `json:"after_count"`
	Message     string `json:"message"`
}


func (m MessageRequest) ToDomain() domain.Message {
	return domain.Message{
		Code:     m.Code,
		Type:     m.Type,
		Category: m.Category,
		Module:   m.Module,
		Title:    m.Title,
		Content:  m.Content,
		Active:   m.Active,
	}
}


func ToMessageResponse(m *domain.Message) MessageResponse {
	return MessageResponse{
		ID:       m.ID,
		Code:     m.Code,
		Type:     m.Type,
		Category: m.Category,
		Module:   m.Module,
		Title:    m.Title,
		Content:  m.Content,
		Active:   m.Active,
	}
}


func ToMessageListResponse(messages []domain.Message) MessageListResponse {
	responses := make([]MessageResponse, len(messages))
	for i, msg := range messages {
		responses[i] = ToMessageResponse(&msg)
	}
	return MessageListResponse{
		Messages: responses,
		Count:    len(responses),
	}
}
