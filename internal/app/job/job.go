package job

import (
	"bytes"
	"context"
	"encoding/gob"
)

type Processor func(ctx context.Context, filePaths []string) error

type Job struct {
	ID        string
	FilePaths []string
	Processor Processor
}

func Encode(job Job) ([]byte, error) {
	var buffer bytes.Buffer
	var encoder = gob.NewEncoder(&buffer)

	if err := encoder.Encode(job); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func Decode(job []byte, result *Job) error {
	var buffer = bytes.NewBuffer(job)
	var decoder = gob.NewDecoder(buffer)
	return decoder.Decode(result)
}
