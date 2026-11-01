package task

import (
	"testing"
)

func TestTaskSerialization(t *testing.T) {
	original, err := NewTask(
		"email:send",
		map[string]string{"to": "bob@example.com"},
	)
	if err != nil {
		t.Fatal(err)
	}

	encoded, err := original.Serialize()
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := Deserialize(encoded)
	if err != nil {
		t.Fatal(err)
	}

	if original.ID != decoded.ID {
		t.Fatal("Task ID changed during serialization")
	}

	if original.Type != decoded.Type {
		t.Fatal("Task type changed during serialization")
	}

	if string(original.Payload) != string(decoded.Payload) {
		t.Fatal("Task payload changed during serialization")
	}

	if !original.CreatedAt.Equal(decoded.CreatedAt) {
		t.Fatal("Task timestamp changed during serialization")
	}
}
