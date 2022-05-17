package bridge

import (
	"context"
	"errors"

	"github.com/parinpan/romusha/definition"
)

var (
	ErrNoBridger = errors.New("assignor has no bridger for member")
)

type bridger interface {
	Assign(ctx context.Context, envelope *definition.JobEnvelope) (resp *definition.Response, err error)
}

type Manager struct {
	bridger map[string]bridger
}

func (m *Manager) AssignByHost(ctx context.Context, host string, envelope *definition.JobEnvelope) (resp *definition.Response, err error) {
	bridger, err := m.getBridger(host)
	if err != nil {
		return nil, ErrNoBridger
	}

	return bridger.Assign(ctx, envelope)
}

func (m *Manager) Add(host string, b bridger) {
	m.bridger[host] = b
}

func (m *Manager) Remove(host string) {
	delete(m.bridger, host)
}

func (m *Manager) getBridger(host string) (bridger, error) {
	if b, ok := m.bridger[host]; ok {
		return b, nil
	}

	return nil, ErrNoBridger
}
