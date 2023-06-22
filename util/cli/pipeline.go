package cli

import (
	"bytes"
	cmdchain "github.com/rainu/go-command-chain"
	"os"
)

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
