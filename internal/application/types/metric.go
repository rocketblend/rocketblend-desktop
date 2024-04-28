package types

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type (
	Metric struct {
		ID         uuid.UUID `json:"id"`
		Domain     string    `json:"domain"`
		Name       string    `json:"name"`
		Value      int       `json:"value"`
		RecordedAt time.Time `json:"recordedAt"`
	}

	GetMetricOpts struct {
		ID uuid.UUID `json:"id"`
	}

	GetMetricResult struct {
		Metric *Metric `json:"metric,omitempty"`
	}

	CreateMetricOpts struct {
		Domain string `json:"domain,omitempty"`
		Name   string `json:"Name,omitempty"`
		Value  int    `json:"value,omitempty"`
	}

	DeleteMetricOpts struct {
		ID uuid.UUID `json:"id"`
	}

	FilterMetricsOpts struct {
		Domain    string    `json:"domain,omitempty"`
		Name      string    `json:"name,omitempty"`
		StartTime time.Time `json:"startTime,omitempty"`
		EndTime   time.Time `json:"endTime,omitempty"`
	}

	ListMetricsResult struct {
		Metrics []*Metric `json:"metrics,omitempty"`
	}

	AggregateMetricsResult struct {
		Aggregate *Aggregate `json:"aggregate,omitempty"`
	}

	Aggregate struct {
		Domain string  `json:"domain"`
		Name   string  `json:"name"`
		Sum    int     `json:"sum"`
		Avg    float64 `json:"avg"`
		Count  int     `json:"count"`
		Min    int     `json:"min"`
		Max    int     `json:"max"`
	}

	Tracker interface {
		GetMetric(ctx context.Context, opts *GetMetricOpts) (*GetMetricResult, error)
		ListMetrics(ctx context.Context, opts *FilterMetricsOpts) (*ListMetricsResult, error)
		AggregateMetrics(ctx context.Context, opts *FilterMetricsOpts) (*AggregateMetricsResult, error)

		CreateMetric(ctx context.Context, opts *CreateMetricOpts) error
		DeleteMetric(ctx context.Context, opts *DeleteMetricOpts) error
	}
)
