package thumbnail

import (
	"bytes"
	cmdchain "github.com/rainu/go-command-chain"
	"os"
)

func CreateImageThumbnail(data []byte) []byte {
	var reader = bytes.NewReader(data)
	var writer = bytes.NewBuffer(nil)

	err := cmdchain.Builder().
		Join("ffmpeg", "-y", "-hide_banner", "-loglevel", "panic",
			"-i", "pipe:0",
			"-vf", "scale=72:72",
			"-f", "image2pipe",
			"pipe:1").
		WithInjections(reader).Finalize().
		WithError(os.Stdout).WithOutput(writer).Run()
	if err != nil {
		return nil
	}

	return writer.Bytes()
}

func CreateVideoThumbnail(data []byte) []byte {
	var reader = bytes.NewReader(data)
	var writer = bytes.NewBuffer(nil)

	err := cmdchain.Builder().
		Join("ffmpeg", "-y", "-hide_banner", "-loglevel", "panic",
			"-f", "mp4", "-i", "pipe:0",
			"-ss", "00:00:00", "-t", "00:00:01",
			"-vf", "'select=gt(scene\\,0.4)'",
			"-frames:v", "5",
			"-fps_mode", "vfr",
			"-vf", "scale=72:72",
			"-f", "image2pipe",
			"pipe:1").
		WithInjections(reader).Finalize().
		WithError(os.Stdout).WithOutput(writer).Run()
	if err != nil {
		return nil
	}

	return writer.Bytes()
}
