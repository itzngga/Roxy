package thumbnail

import (
	"bytes"
	"fmt"
	"github.com/itzngga/roxy/util/cli"
	"github.com/liujiawm/graphics-go/graphics"
	"image"
	"image/jpeg"
)

func CreateImageThumbnail(data []byte) []byte {
	img, _, _ := image.Decode(bytes.NewReader(data))

	dstImage := image.NewRGBA(image.Rect(0, 0, 72, 72))
	err := graphics.Thumbnail(dstImage, img)
	if err != nil {
		fmt.Println(err)
		return data
	}

	result := bytes.Buffer{}
	jpeg.Encode(&result, dstImage, &jpeg.Options{jpeg.DefaultQuality})

	return result.Bytes()
}

func CreateVideoThumbnail(data []byte) []byte {
	dataResult := cli.FfmpegPipeline(data,
		"-y", "-hide_banner", "-loglevel", "panic",
		"-i", "pipe:0",
		"-map_metadata", "-1",
		"-ss", "1",
		"-vframes", "1",
		"-f", "image2",
		"-s", "72x72",
		"pipe:1",
	)

	return dataResult
}
