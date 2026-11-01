package task

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	CreatedAt time.Time       `json:"created_at"`
}

func NewTask(taskType string, payload any) (*Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &Task{
		ID:        uuid.NewString(),
		Type:      taskType,
		Payload:   json.RawMessage(data),
		CreatedAt: time.Now().UTC(),
	}, nil
}

func (t *Task) Serialize() ([]byte, error) {
	return json.Marshal(t)
}

func Deserialize(data []byte) (*Task, error) {
	var t Task

	if err := json.Unmarshal(data, &t); err != nil {
		return nil, err
	}

	return &t, nil
}
