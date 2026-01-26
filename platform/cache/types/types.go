package types

import "context"

type MessageType string

const (
	TypeError   MessageType = "error"
	TypeSuccess MessageType = "success"
	TypeWarning MessageType = "warning"
	TypeInfo    MessageType = "info"
	TypeDebug   MessageType = "debug"
)

type CachedMessage struct {
	ID       string
	Code     string
	Type     MessageType
	Category string
	Module   string
	Title    string
	Content  string
	Active   bool
}

type MessageResponse struct {
	Code    string      `json:"code"`
	Type    MessageType `json:"type"`
	Title   string      `json:"title"`
	Content string      `json:"content"`
}

type MessageCacheRepository interface {
	GetAllActiveForCache(ctx context.Context) ([]CachedMessage, error)
	GetByCodeForCache(ctx context.Context, code string) (*CachedMessage, error)
	GetByCodeWithStatusForCache(ctx context.Context, code string) (*CachedMessage, error)
}
