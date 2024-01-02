package application

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/rocketblend/rocketblend-desktop/internal/application/buffermanager"
)

type EventBufferWriter struct {
	bufferManager buffermanager.BufferManager
}

func BufferWriter(bm buffermanager.BufferManager) io.Writer {
	return &EventBufferWriter{bufferManager: bm}
}

func (cw *EventBufferWriter) Write(p []byte) (n int, err error) {
	var logData map[string]interface{}
	if err := json.Unmarshal(p, &logData); err != nil {
		return 0, err
	}

	message, okMsg := logData["message"].(string)
	level, okLevel := logData["level"].(string)
	timeStr, okTime := logData["time"].(string)
	if !okMsg || !okLevel || !okTime {
		return 0, fmt.Errorf("log message does not contain required 'message' or 'level' or 'time' field")
	}

	delete(logData, "message")
	delete(logData, "level")
	delete(logData, "time")

	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return 0, err
	}

	cw.bufferManager.AddData(LogEvent{
		Level:   level,
		Message: message,
		Time:    parsedTime,
		Fields:  logData,
	})

	return len(p), nil
}
