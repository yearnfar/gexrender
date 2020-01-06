package gexrender

import (
	"encoding/json"
)

// Action 其他操作
type Action struct {
	Action    string          `json:"action" validate:"required"`
	Parameter json.RawMessage `json:"parameter" validate:"required"`
}
