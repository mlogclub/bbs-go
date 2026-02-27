package respath

import "path/filepath"

const (
	ResDirName       = "res"
	UploadsDirName   = "uploads"
	UploadsURLPrefix = "/res/uploads/"
)

func ResDir() string {
	return filepath.Join(".", ResDirName)
}

func UploadsDir() string {
	return filepath.Join(ResDir(), UploadsDirName)
}

func UploadsPath(parts ...string) string {
	paths := []string{UploadsDir()}
	paths = append(paths, parts...)
	return filepath.Join(paths...)
}
