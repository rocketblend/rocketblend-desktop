package types

import "github.com/google/uuid"

const DetailFileName = "rocketdesk.json"

type (
	// ThumbnailSettings struct {
	// 	Width      int    `json:"width,omitempty"`
	// 	Height     int    `json:"height,omitempty"`
	// 	StartFrame int    `json:"startFrame,omitempty"`
	// 	EndFrame   int    `json:"endFrame,omitempty"`
	// 	RenderType string `json:"renderType,omitempty"`
	// }

	Detail struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
		Tags []string  `json:"tags,omitempty"`
		//ThumbnailSettings *ThumbnailSettings `json:"thumbnailSettings,omitempty"`
		ThumbnailPath string `json:"thumbnailPath,omitempty"`
		SplashPath    string `json:"splashPath,omitempty"`
	}
)
