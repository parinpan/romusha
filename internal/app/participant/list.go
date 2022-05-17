package participant

import (
	"context"
	"sync"

	"github.com/parinpan/romusha/definition"
)

var (
	mutex = &sync.Mutex{}
)

type List map[string]struct{}

func (l List) GetAll(_ context.Context) List {
	return l
}

func (l List) Add(_ context.Context, member *definition.Member) (err error) {
	mutex.Lock()
	defer mutex.Unlock()
	l[member.Host] = struct{}{}
	return
}

func (l List) Remove(_ context.Context, member *definition.Member) (err error) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(l, member.Host)
	return
}
