package domain

import (
	"time"

	uuid "github.com/champion19/api-flighthours/tools/utils"
)

type MessageType string

const (
	TypeError   MessageType = "ERROR"
	TypeSuccess MessageType = "EXITO"
	TypeWarning MessageType = "WARNING"
	TypeInfo    MessageType = "INFO"
	TypeDebug   MessageType = "DEBUG"
)

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

func (m *Message) SetID() {
	m.ID = uuid.Generate()
}

func (m *Message) ToLogger() []string {
	return []string{
		"id:" + m.ID,
		"code:" + m.Code,
		"type:" + string(m.Type),
		"module:" + m.Module,
	}
}

func (m *Message) Validate() error {
	if m.Code == "" {
		return ErrMessageCodeRequired
	}
	return nil
}
