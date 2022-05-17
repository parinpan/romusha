package definition

import (
	"context"
)

type Executor interface {
	Execute(ctx context.Context, source string)
}
