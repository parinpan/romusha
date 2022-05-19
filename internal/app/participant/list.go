package participant

import (
	"context"
	"sync"

	"github.com/parinpan/romusha/definition"
)

var (
	mutex = &sync.Mutex{}
)

type List map[string]definition.Status

func (l List) GetAll(_ context.Context) List {
	return l
}

func (l List) HostsByStatus(_ context.Context, status definition.Status) (hosts []string) {
	for host, s := range l {
		if s.String() == status.String() {
			hosts = append(hosts, host)
		}
	}

	return
}

func (l List) Add(_ context.Context, member *definition.Member, status definition.Status) (err error) {
	mutex.Lock()
	defer mutex.Unlock()
	l[member.Host] = status
	return
}

func (l List) Remove(_ context.Context, member *definition.Member) (err error) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(l, member.Host)
	return
}
