package application

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/listoption"
)

var validExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".bmp":  true,
	".svg":  true,
	".webp": true,
}

type FileLoader struct {
	logger logger.Logger

	store searchstore.Store
}

func NewFileLoader(logger logger.Logger, store searchstore.Store) (*FileLoader, error) {
	return &FileLoader{
		logger: logger,
		store:  store,
	}, nil
}

func (h *FileLoader) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		h.respondWithError(res, http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		return
	}

	path := strings.TrimPrefix(req.URL.Path, "/")
	pathParts := strings.Split(path, "/")

	if len(pathParts) != 2 || pathParts[0] != "system" {
		h.respondWithError(res, http.StatusBadRequest, "Invalid path", nil)
		return
	}

	resourcePath := pathParts[1]
	result, err := h.store.List(listoption.WithResource(resourcePath), listoption.WithSize(1))
	if err != nil {
		h.respondWithError(res, http.StatusInternalServerError, "Could not list resources", err)
		return
	}

	if len(result) == 0 {
		h.respondWithError(res, http.StatusNotFound, "Resource not found", nil)
		return
	}

	if !isValidWebImage(resourcePath) {
		h.respondWithError(res, http.StatusBadRequest, "Invalid resource file type", fmt.Errorf("invalid file path: %s", resourcePath))
		return
	}

	h.logger.Debug("Requested file:", map[string]interface{}{
		"hitIndexes":   len(result),
		"resourcePath": resourcePath,
		"method":       req.Method,
		"remoteAddr":   req.RemoteAddr,
		"userAgent":    req.UserAgent(),
	})

	http.ServeFile(res, req, resourcePath)
}

func (h *FileLoader) respondWithError(w http.ResponseWriter, statusCode int, message string, err error) {
	if err != nil {
		h.logger.Error(message, map[string]interface{}{
			"error": err.Error(),
		})
	}

	w.WriteHeader(statusCode)
	w.Write([]byte(fmt.Sprintf("%s: %v", message, err)))
}

func isValidWebImage(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	return validExtensions[ext]
}
