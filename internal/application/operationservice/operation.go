package operationservice

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/indextype"
)

type (
	Operation struct {
		ID        uuid.UUID   `json:"id"`
		Completed bool        `json:"completed"`
		ErrorMsg  string      `json:"error,omitempty"`
		Result    interface{} `json:"result,omitempty"`
	}
)

func (o Operation) ToSearchIndex() (*searchstore.Index, error) {
	data, err := json.Marshal(o)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal OperationStatus: %w", err)
	}

	state := 0
	if o.Completed {
		state = 1
	}

	return &searchstore.Index{
		ID:    o.ID,
		Type:  indextype.Operation,
		State: state,
		Error: o.ErrorMsg,
		Data:  string(data),
	}, nil
}
