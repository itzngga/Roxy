package cli

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

func ExecPipeline(command string, data []byte, arg ...string) []byte {
	cmd := exec.Command(command, arg...)

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

	stderr, err := cmd.StderrPipe()
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

	errBytes := make([]byte, 0)

	defer stderr.Close()
	errBytes, err = io.ReadAll(stderr)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	if errBytes != nil {
		fmt.Println(string(errBytes))
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return outBytes
}
