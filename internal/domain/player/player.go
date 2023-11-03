package player

import (
	"github.com/toledoom/gork/pkg/entity"
	"github.com/toledoom/gork/pkg/event"
)

type Player struct {
	ag *entity.Aggregate

	ID    string
	Name  string
	Score int64
}

func (p *Player) AddEvent(e event.Event) {
	p.ag.AddEvent(e)
}

func (p *Player) GetEvents() []event.Event {
	return p.ag.GetEvents()
}

func New(id, name string) *Player {
	return &Player{
		ag: &entity.Aggregate{},

		ID:   id,
		Name: name,
	}
}

var _ entity.Entity = (*Player)(nil)

type Repository interface {
	Add(p *Player) error
	Update(p *Player) error
	Delete(p *Player) error
	GetByID(id string) (*Player, error)
}
