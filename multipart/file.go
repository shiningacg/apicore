package multipart

import (
	"io"
	"os"
)

type File struct {
	name string
	io.ReadCloser
}

func (f *File) FileName() string {
	return f.name
}

func (f *File) Close() error {
	return f.ReadCloser.Close()
}

type IOFile struct {
	name string
	*os.File
}

func (f *IOFile) Close() error {
	var err error
	err = f.File.Close()
	err = os.Remove(f.name)
	return err
}

type BufferFile struct {
	io.Reader
}

func (b *BufferFile) Close() error {
	return nil
}
