package cli

import (
	"bytes"
	cmdchain "github.com/rainu/go-command-chain"
	"os"
)

// ExecPipeline execute command with stdin input and capture stdout to byte
func ExecPipeline(cmd string, data []byte, params ...string) ([]byte, error) {
	var reader = bytes.NewReader(data)
	var writer = bytes.NewBuffer(nil)

	err := cmdchain.Builder().
		Join(cmd, params...).
		WithInjections(reader).Finalize().
		WithError(os.Stdout).WithOutput(writer).Run()
	if err != nil {
		return nil, err
	}

	return writer.Bytes(), nil
}

// Exec execute command without input data and capture stdout to byte
func Exec(cmd string, params ...string) ([]byte, error) {
	var writer = bytes.NewBuffer(nil)

	err := cmdchain.Builder().
		Join(cmd, params...).Finalize().
		WithError(os.Stdout).WithOutput(writer).Run()
	if err != nil {
		return nil, err
	}

	return writer.Bytes(), nil
}
