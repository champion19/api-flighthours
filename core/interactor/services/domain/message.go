package domain

import (
	"time"

	uuid "github.com/champion19/flighthours-api/tools/utils"
)

	// MessageType represents system message types
type MessageType string

const (
	TypeError   MessageType = "ERROR"
	TypeSuccess MessageType = "EXITO"
	TypeWarning MessageType = "WARNING"
	TypeInfo    MessageType = "INFO"
	TypeDebug   MessageType = "DEBUG"
)

// Message represents a system message in the domain
type Message struct {
	ID        string      `json:"id"`
	Code      string      `json:"code"`
	Type      MessageType `json:"type"`
	Category  string      `json:"category"`
	Module    string      `json:"module"`
	Title     string      `json:"title"`
	Content   string      `json:"content"`
	Active    bool        `json:"active"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// SetID generates a new UUID for the message
func (m *Message) SetID() {
	m.ID = uuid.Generate()
}

// ToLogger returns a slice of strings for logging
func (m *Message) ToLogger() []string {
	return []string{
		"id:" + m.ID,
		"code:" + m.Code,
		"type:" + string(m.Type),
		"module:" + m.Module,
	}
}

// TODO: quitar estas validaciones de aqui
// IsValid validates the message fields
func (m *Message) Validate() error {
	if m.Code == "" {
		return ErrMessageCodeRequired
	}
	return nil
}
