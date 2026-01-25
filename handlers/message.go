package handlers
import (

	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
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

// MessageResponse represents the response payload for a message
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

// MessageListResponse represents the response for listing messages
type MessageListResponse struct {
	Messages []MessageResponse `json:"messages"`
	Count    int               `json:"count"`
	Links    []Link            `json:"_links,omitempty"`
}

// MessageCreatedResponse represents the response for message creation
type MessageCreatedResponse struct {
	ID    string `json:"id"`
	Links []Link `json:"_links"`

}

// MessageUpdatedResponse represents the response for message update
type MessageUpdatedResponse struct {
	Links []Link `json:"_links"`

}

// MessageDeletedResponse represents the response for message deletion
type MessageDeletedResponse struct {
	// Empty struct - message comes from unified messaging system
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

// ToResponse converts domain.Message to MessageResponse
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

// ToMessageListResponse converts slice of domain.Message to MessageListResponse
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
