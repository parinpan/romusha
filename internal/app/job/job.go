package job

import (
	"bytes"
	"context"
	"encoding/gob"
)

type Processor func(ctx context.Context, filePaths []string) error

func Encode[T](data *T) ([]byte, error) {
	var buffer bytes.Buffer
	var encoder = gob.NewEncoder(&buffer)

	if err := encoder.Encode(data); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func Decode[T](job []byte, result *T) error {
	var buffer = bytes.NewBuffer(job)
	var decoder = gob.NewDecoder(buffer)
	return decoder.Decode(result)
}
