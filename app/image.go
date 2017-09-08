package app

import (
	"gopkg.in/gographics/imagick.v2/imagick"
)

var ImgCli *ImageClient

func GetImgClient() *ImageClient {
	return ImgCli
}

type ImageClient struct {
	mw *imagick.MagickWand
}

func InitImageClient() *ImageClient {
	mw := imagick.NewMagickWand()
	return &ImageClient{
		mw: mw,
	}
}

func TerminateImageMagick() {
	imagick.Terminate()
}

func (self *ImageClient) ReadImage(imagePath string) error {
	err := self.mw.ReadImage(imagePath)
	if err != nil {
		return err
	}
	// Get original logo size
	width := self.mw.GetImageWidth()
	height := self.mw.GetImageHeight()

	// Calculate half the size
	hWidth := uint(width / 4)
	hHeight := uint(height / 4)

	// Resize the image using the Lanczos filter
	// The blur factor is a float, where > 1 is blurry, < 1 is sharp
	err = self.mw.ResizeImage(hWidth, hHeight, imagick.FILTER_LANCZOS, 1)

	if err != nil {
		return err
	}

	err = self.mw.SetImageCompressionQuality(80)
	if err != nil {
		return err
	}

	self.mw.WriteImage("/tmp/test.jpg")

	return nil
}
