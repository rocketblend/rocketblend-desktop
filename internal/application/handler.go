package application

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend-desktop/internal/application/factory"
	"github.com/rocketblend/rocketblend-desktop/internal/application/searchstore/listoption"
)

const DynamicResourcePath = "system/"

var validExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".bmp":  true,
	".svg":  true,
	".webp": true,
	".webm": true,
}

type FileLoader struct {
	cacheTimeout int

	logger  logger.Logger
	factory factory.Factory
}

func NewFileLoader(logger logger.Logger, factory factory.Factory, cacheTimeout int) (*FileLoader, error) {
	return &FileLoader{
		logger:       logger,
		factory:      factory,
		cacheTimeout: cacheTimeout,
	}, nil
}

func (h *FileLoader) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		h.respondWithError(res, http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		return
	}

	path := strings.TrimPrefix(req.URL.Path, "/")
	if !strings.HasPrefix(path, DynamicResourcePath) {
		h.respondWithError(res, http.StatusBadRequest, "Invalid path", nil)
		return
	}

	path = filepath.ToSlash(strings.TrimPrefix(path, DynamicResourcePath))
	if !isValidWebImage(path) {
		h.respondWithError(res, http.StatusBadRequest, "Invalid resource file type", fmt.Errorf("invalid file type: %s", path))
		return
	}

	h.logger.Debug("Requested file:", map[string]interface{}{
		"resource":   path,
		"method":     req.Method,
		"remoteAddr": req.RemoteAddr,
		"userAgent":  req.UserAgent(),
	})

	store, err := h.factory.GetSearchStore()
	if err != nil {
		h.respondWithError(res, http.StatusInternalServerError, "Could not get search store", err)
		return
	}

	result, err := store.List(listoption.WithResource(path), listoption.WithSize(1))
	if err != nil {
		h.respondWithError(res, http.StatusInternalServerError, "Could not list resources", err)
		return
	}

	if len(result) == 0 {
		h.respondWithError(res, http.StatusNotFound, "Resource not found", nil)
		return
	}

	res.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", h.cacheTimeout))

	http.ServeFile(res, req, path)
}

func (h *FileLoader) respondWithError(w http.ResponseWriter, statusCode int, message string, err error) {
	if err != nil {
		h.logger.Error(message, map[string]interface{}{
			"error": err.Error(),
		})
	} else {
		h.logger.Warn(message)
	}

	w.WriteHeader(statusCode)
	w.Write([]byte(fmt.Sprintf("%s: %v", message, err)))
}

func isValidWebImage(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	return validExtensions[ext]
}
