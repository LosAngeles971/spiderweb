package knowledge

import (
	"errors"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	//"sync"
)

type Knowledge struct {
	Events 	[]*Event `yaml:"knowlegde"`
	//lock  sync.RWMutex
}

func (k *Knowledge) Load(kfile string) error {
	b, err := ioutil.ReadFile(kfile)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(b, k)
	if err != nil {
		return err
	}
	for _, e := range k.Events {
		for _, r := range e.Effects {
			ef, ok := k.GetEvent(r.Indirect)
			if !ok {
				return errors.New("Missing event: " + r.Indirect)
			}
			r.Effect = ef
		}
	}
	return nil
}

// This function search for a Event
func (u Knowledge) GetEvent(id string) (*Event, bool) {
	for _, e := range u.Events {
		if e.GetID() == id {
			return e, true
		}
	}
	return nil, false
}

// This function search for a Event
/*
func (u Knowledge) GetRelationship(cause Event, effect Event) (Relationship, bool) {
	for _, r := range u.Relationships {
		if r.Cause.GetID() == cause.GetID() && r.Effect.GetID() == effect.GetID() {
			return r, true
		}
	}
	return Relationship{}, false
}
*/

// This function returns all events that cause the given effect
func (u Knowledge) IsEffectOf(effect *Event) ([]*Event) {
	result := []*Event{}
	for _, cs := range u.Events {
		for _, r := range cs.Effects {
			if r.Effect.GetID() == effect.GetID() && r.GetWeight() > 0.0 {
				result = append(result, cs)
			}
		}
	}
	return result
}