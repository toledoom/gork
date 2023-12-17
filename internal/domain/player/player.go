package player

import (
	"github.com/toledoom/gork/pkg/gork"
)

type Player struct {
	ag *gork.Aggregate

	ID    string
	Name  string
	Score int64
}

func (p *Player) AddEvent(e gork.Event) {
	p.ag.AddEvent(e)
}

func (p *Player) GetEvents() []gork.Event {
	return p.ag.GetEvents()
}

func New(id, name string) *Player {
	return &Player{
		ag: &gork.Aggregate{},

		ID:   id,
		Name: name,
	}
}

var _ gork.Entity = (*Player)(nil)

type Repository interface {
	Add(p *Player) error
	Update(p *Player) error
	Delete(p *Player) error
	GetByID(id string) (*Player, error)
}
