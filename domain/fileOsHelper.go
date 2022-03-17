package domain

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/nfnt/resize"
)

type IFileOsHelper interface {
	DecodeImage(io.Reader) (image.Image, string, error)
	ResizeAndSaveFileToLocale(string, image.Image, string, string) (string, error)
	MkDirAll(string, fs.FileMode) error
	Create(string) (*os.File, error)
	Copy(io.Writer, io.Reader) (int64, error)
}

type FileOsHelper struct {
}

func NewFileOsHelper() *FileOsHelper {
	return &FileOsHelper{}
}

func (fos *FileOsHelper) DecodeImage(r io.Reader) (image.Image, string, error) {
	return image.Decode(r)
}

func (fos *FileOsHelper) ResizeAndSaveFileToLocale(size string, originalProfilePicture image.Image, userId string, fileType string) (string, error) {
	var resizedProfilePicture image.Image
	fileExtension := strings.Split(fileType, "/")[1]
	resizedProfilePictureUrl := "./profile_pictures/" + size + "-profile-picture-" + userId + "." + fileExtension
	switch size {
	case "small":
		resizedProfilePicture = resize.Resize(150, 150, originalProfilePicture, resize.Lanczos3)
	case "average":
		resizedProfilePicture = resize.Resize(400, 400, originalProfilePicture, resize.Lanczos3)
	case "large":
		resizedProfilePicture = resize.Resize(800, 800, originalProfilePicture, resize.Lanczos3)
	}

	resizedProfilePictureFile, err := os.Create(resizedProfilePictureUrl)
	if err != nil {
		return resizedProfilePictureUrl, err
	}
	if fileExtension == "jpeg" {
		jpeg.Encode(resizedProfilePictureFile, resizedProfilePicture, nil)
	}
	if fileExtension == "png" {
		png.Encode(resizedProfilePictureFile, resizedProfilePicture)
	}
	if fileExtension == "gif" {
		gif.Encode(resizedProfilePictureFile, resizedProfilePicture, nil)
	}

	resizedProfilePictureFile.Close()
	return resizedProfilePictureUrl, nil
}

func (fos *FileOsHelper) MkDirAll(path string, perm fs.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (fos *FileOsHelper) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (fos *FileOsHelper) Copy(dst io.Writer, src io.Reader) (int64, error) {
	return io.Copy(dst, src)
}
