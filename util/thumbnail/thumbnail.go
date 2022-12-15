package thumbnail

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/liujiawm/graphics-go/graphics"
	"image"
	"image/jpeg"
	"io"
	"os/exec"
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
	cmd := exec.Command("ffmpeg", "-y",
		"-hide_banner", "-loglevel", "panic",
		"-i", "pipe:0",
		"-map_metadata", "-1",
		"-ss", "1",
		"-vframes", "1",
		"-f", "image2",
		"-s", "72x72",
		"pipe:1",
	)

	out, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	in, err := cmd.StdinPipe()
	writer := bufio.NewWriter(in)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	go func() {
		defer writer.Flush()
		defer in.Close()
		_, err = writer.Write(data)
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	outBytes := make([]byte, 0)

	defer out.Close()
	outBytes, err = io.ReadAll(out)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return outBytes
}
