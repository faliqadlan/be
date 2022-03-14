package s3

import "mime/multipart"

type TaskS3M interface {
	UploadFileToS3(fileHeader multipart.FileHeader) (string, error)
	UpdateFileS3(name string, fileHeader multipart.FileHeader) string
	DeleteFileS3(name string) string
}
