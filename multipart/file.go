package multipart

import "mime/multipart"

type File struct {
	*multipart.Part
}
