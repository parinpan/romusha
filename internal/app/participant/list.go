package participant

import (
	"context"
	"sync"
)

var (
	mutex = &sync.Mutex{}
)

type Host string
type Endpoint string
type List map[Host]Endpoint

func (l List) GetAll(_ context.Context) List {
	return l
}

func (l List) GetEndpoint(_ context.Context, host Host) Endpoint {
	mutex.Lock()
	defer mutex.Unlock()
	return l[host]
}

func (l List) Add(_ context.Context, member Member) (err error) {
	mutex.Lock()
	defer mutex.Unlock()
	l[Host(member.Host)] = Endpoint(member.Endpoint)
	return
}

func (l List) Remove(_ context.Context, member Member) (err error) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(l, Host(member.Host))
	return
}
