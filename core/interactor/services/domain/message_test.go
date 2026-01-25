package domain

import (
	"testing"
)

func TestMessage_SetID(t *testing.T) {
	msg := &Message{}

	if msg.ID != "" {
		t.Fatal("expected ID to be empty initially")
	}

	msg.SetID()

	if msg.ID == "" {
		t.Fatal("expected ID to be set after SetID()")
	}
}

func TestMessage_ToLogger(t *testing.T) {
	msg := &Message{
		ID:     "msg-123",
		Code:   "TEST_001",
		Type:   TypeError,
		Module: "auth",
	}

	result := msg.ToLogger()

	if len(result) != 4 {
		t.Fatalf("expected 4 elements, got %d", len(result))
	}

	expectedID := "id:msg-123"
	if result[0] != expectedID {
		t.Fatalf("expected %s, got %s", expectedID, result[0])
	}

	expectedCode := "code:TEST_001"
	if result[1] != expectedCode {
		t.Fatalf("expected %s, got %s", expectedCode, result[1])
	}

	expectedType := "type:ERROR"
	if result[2] != expectedType {
		t.Fatalf("expected %s, got %s", expectedType, result[2])
	}

	expectedModule := "module:auth"
	if result[3] != expectedModule {
		t.Fatalf("expected %s, got %s", expectedModule, result[3])
	}
}

func TestMessage_ToLogger_EmptyFields(t *testing.T) {
	msg := &Message{}

	result := msg.ToLogger()

	if len(result) != 4 {
		t.Fatalf("expected 4 elements, got %d", len(result))
	}

	if result[0] != "id:" {
		t.Fatalf("expected 'id:', got %s", result[0])
	}
	if result[1] != "code:" {
		t.Fatalf("expected 'code:', got %s", result[1])
	}
	if result[2] != "type:" {
		t.Fatalf("expected 'type:', got %s", result[2])
	}
	if result[3] != "module:" {
		t.Fatalf("expected 'module:', got %s", result[3])
	}
}

func TestMessage_Validate(t *testing.T) {
	t.Run("empty code returns ErrMessageCodeRequired", func(t *testing.T) {
		msg := &Message{Code: ""}
		err := msg.Validate()

		if err != ErrMessageCodeRequired {
			t.Fatalf("expected ErrMessageCodeRequired, got %v", err)
		}
	})

	t.Run("valid code returns nil", func(t *testing.T) {
		msg := &Message{Code: "TEST_001"}
		err := msg.Validate()

		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	})
}

func TestMessageType_Constants(t *testing.T) {
	tests := []struct {
		name     string
		msgType  MessageType
		expected string
	}{
		{"TypeError", TypeError, "ERROR"},
		{"TypeSuccess", TypeSuccess, "EXITO"},
		{"TypeWarning", TypeWarning, "WARNING"},
		{"TypeInfo", TypeInfo, "INFO"},
		{"TypeDebug", TypeDebug, "DEBUG"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.msgType) != tt.expected {
				t.Fatalf("expected %s, got %s", tt.expected, string(tt.msgType))
			}
		})
	}
}
