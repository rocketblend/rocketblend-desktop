package fileserver

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/rocketblend/rocketblend-desktop/internal/application/store/listoption"
	"github.com/rocketblend/rocketblend-desktop/internal/application/types"
)

const DynamicResourcePath = "system/"

var ValidExtensions = map[string]struct{}{
	".jpg":  {},
	".jpeg": {},
	".png":  {},
	".gif":  {},
	".bmp":  {},
	".svg":  {},
	".webp": {},
	".webm": {},
}

type (
	dependecies struct {
		logger types.Logger
		store  types.Store
	}

	Options struct {
		Container    types.Container
		CacheTimeout int
	}

	Option func(*Options)

	Handler struct {
		logger       types.Logger
		store        types.Store
		cacheTimeout int
	}
)

func WithContainer(container types.Container) Option {
	return func(o *Options) {
		o.Container = container
	}
}

func WithCacheTimeout(timeout int) Option {
	return func(o *Options) {
		o.CacheTimeout = timeout
	}
}

func New(opts ...Option) (*Handler, error) {
	options := &Options{
		CacheTimeout: 3600, // 1 hour
	}

	for _, opt := range opts {
		opt(options)
	}

	if options.Container == nil {
		return nil, errors.New("container is required")
	}

	dependencies, err := setupDependencies(options.Container)
	if err != nil {
		return nil, fmt.Errorf("failed to setup dependencies: %w", err)
	}

	return &Handler{
		logger:       dependencies.logger,
		store:        dependencies.store,
		cacheTimeout: options.CacheTimeout,
	}, nil
}

func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		h.respondWithError(res, http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		return
	}

	path := strings.TrimPrefix(req.URL.Path, "/")
	if !strings.HasPrefix(path, DynamicResourcePath) {
		h.respondWithError(res, http.StatusBadRequest, "Invalid path", nil)
		return
	}

	path = ensureAbsolutePath(filepath.ToSlash(strings.TrimPrefix(path, DynamicResourcePath)))
	if !isValidFile(path) {
		h.respondWithError(res, http.StatusBadRequest, "Invalid resource file type", fmt.Errorf("invalid file type: %s", path))
		return
	}

	h.logger.Debug("requested file:", map[string]interface{}{
		"resource":   path,
		"method":     req.Method,
		"remoteAddr": req.RemoteAddr,
		"userAgent":  req.UserAgent(),
	})

	result, err := h.store.List(req.Context(), listoption.WithResource(path), listoption.WithSize(1))
	if err != nil {
		h.respondWithError(res, http.StatusInternalServerError, "Could not list resources", err)
		return
	}

	if len(result) == 0 {
		h.respondWithError(res, http.StatusNotFound, "Resource not found", nil)
		return
	}

	h.logger.Debug("serving file:", map[string]interface{}{
		"resource": path,
	})

	res.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", h.cacheTimeout))

	http.ServeFile(res, req, path)
}

func (h *Handler) respondWithError(w http.ResponseWriter, statusCode int, message string, err error) {
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

func isValidFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	_, exists := ValidExtensions[ext]
	return exists
}

// ensureAbsolutePath makes sure the path starts with a '/'.
func ensureAbsolutePath(path string) string {
	if len(path) > 1 && path[1] == ':' {
		return path
	}

	if !strings.HasPrefix(path, "/") {
		return "/" + path
	}

	return path
}

func setupDependencies(container types.Container) (*dependecies, error) {
	logger, err := container.GetLogger()
	if err != nil {
		return nil, err
	}

	store, err := container.GetStore()
	if err != nil {
		return nil, err
	}

	return &dependecies{
		logger: logger,
		store:  store,
	}, nil
}
