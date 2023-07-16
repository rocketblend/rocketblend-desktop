package project

type (
	ThumbnailSettings struct {
		Width      int    `json:"width"`
		Height     int    `json:"height"`
		StartFrame int    `json:"startFrame"`
		EndFrame   int    `json:"endFrame"`
		RenderType string `json:"renderType"`
	}

	Settings struct {
		Name              string             `json:"name"`
		Tags              []string           `json:"tags"`
		ThumbnailSettings *ThumbnailSettings `json:"thumbnailSettings"`
		ThumbnailFilePath string             `json:"thumbnailFilePath"`
	}
)
