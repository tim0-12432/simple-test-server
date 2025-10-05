package web

// UploadResult contains metadata about a saved uploaded file.
type UploadResult struct {
	LocalPath   string `json:"local_path"`
	SafeName    string `json:"safe_name"`
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
}
