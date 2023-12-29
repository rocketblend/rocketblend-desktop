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
		ID        uuid.UUID
		Completed bool
		ErrorMsg  string
		Result    interface{}
	}
)

func (o Operation) ToSearchIndex() (*searchstore.Index, error) {
	data, err := json.Marshal(o)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal OperationStatus: %w", err)
	}

	return &searchstore.Index{
		ID:    o.ID,
		Type:  indextype.Operation,
		Ready: !o.Completed,
		Error: o.ErrorMsg,
		Data:  string(data),
	}, nil
}
