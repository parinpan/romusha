package job

import (
	"bytes"
	"context"
	"encoding/gob"
)

type Processor func(ctx context.Context, filePaths []string) error

func Decode(job []byte, result interface{}) error {
	var buffer = bytes.NewBuffer(job)
	var decoder = gob.NewDecoder(buffer)
	return decoder.Decode(result)
}
