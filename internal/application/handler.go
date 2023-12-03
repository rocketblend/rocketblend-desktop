package application

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/projectservice"
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
	logger         logger.Logger
	projectService projectservice.Service
}

func NewFileLoader(logger logger.Logger, projectService projectservice.Service) (*FileLoader, error) {
	return &FileLoader{
		logger:         logger,
		projectService: projectService,
	}, nil
}

func (h *FileLoader) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		h.respondWithError(res, http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		return
	}

	path := strings.TrimPrefix(req.URL.Path, "/")
	pathParts := strings.Split(path, "/")

	if len(pathParts) != 2 || pathParts[0] != "image" {
		h.respondWithError(res, http.StatusBadRequest, "Invalid path", nil)
		return
	}

	id, err := uuid.Parse(pathParts[1])
	if err != nil {
		h.respondWithError(res, http.StatusBadRequest, "Could not parse project id", err)
		return
	}

	result, err := h.projectService.Get(req.Context(), id)
	if err != nil {
		h.respondWithError(res, http.StatusNotFound, "Could not find project", err)
		return
	}

	if result == nil || result.Project == nil {
		h.respondWithError(res, http.StatusNotFound, "No project found", nil)
		return
	}

	if !isValidWebImage(result.Project.ImagePath) {
		h.respondWithError(res, http.StatusBadRequest, "Invalid image file", fmt.Errorf("invalid file path: %s", result.Project.ImagePath))
		return
	}

	absolutePath := filepath.Join(result.Project.Path, result.Project.ImagePath)
	h.logger.Debug("Requested file:", map[string]interface{}{
		"id":           id,
		"relativePath": result.Project.ImagePath,
		"absolutePath": absolutePath,
		"method":       req.Method,
		"remoteAddr":   req.RemoteAddr,
		"userAgent":    req.UserAgent(),
	})

	http.ServeFile(res, req, absolutePath)
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
