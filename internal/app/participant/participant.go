package participant

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
)

const (
	stateChannel = "participant:state"
)

type pubSubClient interface {
	Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd
	Subscribe(ctx context.Context, channels ...string) *redis.PubSub
}

type subscriberClient interface {
	ReceiveMessage(ctx context.Context) (*redis.Message, error)
}

type Participant struct {
	list List
	psc  pubSubClient
	ssc  func(ctx context.Context) subscriberClient
}

func NewParticipant(redis pubSubClient) *Participant {
	return &Participant{
		list: List{},
		psc:  redis,
		ssc:  subscribe(redis, stateChannel),
	}
}

func (p *Participant) Notify(ctx context.Context, state StateBody) error {
	return p.psc.Publish(ctx, stateChannel, state).Err()
}

func (p *Participant) List(ctx context.Context) List {
	return p.list.GetAll(ctx)
}

func (p *Participant) Watch(ctx context.Context, watchers ...Watcher) (err error) {
	var data StateBody

	for {
		msg, err := p.ssc(ctx).ReceiveMessage(ctx)
		if err != nil {
			break
		}

		if msg == nil {
			continue
		}

		if err := json.Unmarshal([]byte(msg.Payload), &data); err != nil {
			continue
		}

		for _, watcher := range watchers {
			if err := watcher(ctx, data); err != nil {
				continue
			}
		}
	}

	return
}

func subscribe(client pubSubClient, channel ...string) func(ctx context.Context) subscriberClient {
	return func(ctx context.Context) subscriberClient {
		return client.Subscribe(ctx, channel...)
	}
}
