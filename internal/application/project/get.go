package project

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/rocketblend/rocketblend-desktop/internal/application/fileserver"
	"github.com/rocketblend/rocketblend-desktop/internal/application/store/indextype"
	"github.com/rocketblend/rocketblend-desktop/internal/application/types"
)

var validMediaExtensions = map[string]struct{}{
	".jpg":  {},
	".jpeg": {},
	".png":  {},
	".gif":  {},
	".bmp":  {},
	".svg":  {},
	".webp": {},
	".webm": {},
}

var splashKeyWord = "splash"
var thumbnailKeyWord = "thumbnail"

func (r *Repository) GetProject(ctx context.Context, opts *types.GetProjectOpts) (*types.GetProjectResponse, error) {
	result, err := r.get(ctx, opts.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	return &types.GetProjectResponse{
		Project: result,
	}, nil
}

func (r *Repository) get(ctx context.Context, id uuid.UUID) (*types.Project, error) {
	index, err := r.store.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	project, err := convertFromIndex(index)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func convertFromIndex(index *types.Index) (*types.Project, error) {
	var result types.Project
	if err := json.Unmarshal([]byte(index.Data), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func convertToIndex(project *types.Project) (*types.Index, error) {
	data, err := json.Marshal(project)
	if err != nil {
		return nil, err
	}

	// TODO: need a better way to register resources for the http server.
	resources := make([]string, 0, len(project.Media))
	for _, m := range project.Media {
		resources = append(resources, filepath.ToSlash(m.FilePath))
	}

	return &types.Index{
		ID:        project.ID,
		Name:      project.Name,
		Type:      indextype.Project,
		Reference: path.Clean(project.Path),
		Resources: resources,
		Data:      string(data),
	}, nil
}

func findMediaFiles(path string) ([]*types.Media, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}

	var splash bool
	var thumbnail bool

	var files []*types.Media
	visit := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil // skip directories
		}

		ext := strings.ToLower(filepath.Ext(path))
		if _, found := validMediaExtensions[ext]; found {
			media, err := loadMedia(path)
			if err != nil {
				return err
			}

			if media.Splash {
				splash = true
			}

			if media.Thumbnail {
				thumbnail = true
			}

			files = append(files, media)
		}

		return nil
	}

	err := filepath.WalkDir(path, visit)
	if err != nil {
		return nil, err
	}

	if !splash && len(files) > 0 {
		files[0].Splash = true
	}

	if !thumbnail && len(files) > 0 {
		files[0].Thumbnail = true
	}

	return files, nil
}

func loadMedia(path string) (*types.Media, error) {
	return &types.Media{
		FilePath:  path,
		URL:       "/" + filepath.ToSlash(filepath.Join(fileserver.DynamicResourcePath, path)),
		Splash:    containsWordInFilename(path, splashKeyWord),
		Thumbnail: containsWordInFilename(path, thumbnailKeyWord),
	}, nil
}

func containsWordInFilename(path string, word string) bool {
	filename := filepath.Base(path)
	lowerFilename := strings.ToLower(filename)
	lowerWord := strings.ToLower(word)
	return strings.Contains(lowerFilename, lowerWord)
}
