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

	reader = nil
	writer = nil
	return writer.Bytes()
}

func CreateVideoThumbnail(data []byte) []byte {
	var reader = bytes.NewReader(data)
	var writer = bytes.NewBuffer(nil)

	err := cmdchain.Builder().
		Join("ffmpeg", "-y", "-hide_banner", "-loglevel", "panic",
			"-f", "mp4", "-i", "pipe:0",
			"-filter_complex", "'scale=72:72,select=between(t\\,10\\,20)*eq(pict_type\\,I)'",
			"-update", "true",
			"-fps_mode", "vfr",
			"-f", "image2pipe",
			"pipe:1").
		WithInjections(reader).Finalize().
		WithError(os.Stdout).WithOutput(writer).Run()
	if err != nil {
		return nil
	}

	reader = nil
	writer = nil
	return writer.Bytes()
}
