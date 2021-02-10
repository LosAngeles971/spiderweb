package knowledge

import (
	"errors"
	"strings"
	"github.com/Knetic/govaluate"
	ctx "it/losangeles971/joshua/internal/context"
	"it/losangeles971/joshua/internal/math"
)

const (
	EVENT_OUTCOME_TRUE 		= "true"
	EVENT_OUTCOME_FALSE 	= "false"
	EVENT_OUTCOME_UNKNOWN 	= "missing data"
	EVENT_OUTCOME_ERROR 	= "error"
)

/*
A Event or Event is just a Event or event :-)
So it is en established truth, and it can be defined by different charactericts at the same time:

- a name
- a list of equations
- a list of attributes

From a cause/effect perspective, a Event may be a casue and/or an effect.

From model perspective a Event is the vertex of directed graph, where the edges mean cause/effect relationships.
*/

type Event struct {
	ID			string 		`yaml:"id"`
	Statements	[]string 	`yaml:"statements"`
	Assignments	[]string 	`yaml:"assignments"`
	Dependecies	[]Event
}

func parse(a string) (string, *govaluate.EvaluableExpression, error) {
	parts := strings.Split(a, "=")
	if len(parts) != 2 {
		return "", nil, errors.New("Malformed assignment: " + a)
	}
	e, err := govaluate.NewEvaluableExpressionWithFunctions(parts[1], math.Functions)
	if err != nil {
		return "", nil, err
	}
	return strings.TrimSpace(parts[0]), e, nil
}

func (f *Event) Verify(data *ctx.State) (string, error) {
	for _, s := range f.Statements {
		e, err := govaluate.NewEvaluableExpressionWithFunctions(s, math.Functions)
		if err != nil {
			return EVENT_OUTCOME_ERROR, err
		}
		// Check if the function requires variables not present into the context
		for _, v := range e.Vars() {
			if ok := data.Contains(v); !ok {
				return EVENT_OUTCOME_UNKNOWN, nil
			}
		}
		r, err := e.Evaluate(data.State())
		switch r.(type) {
		case bool:
			if !r.(bool) {
				return EVENT_OUTCOME_FALSE, nil
			}
		default:
			return EVENT_OUTCOME_ERROR, errors.New("Statement must be boolean: " + s)
		}
	}
	for _, a := range f.Assignments {
		v, e, err := parse(a)
		if err != nil {
			return EVENT_OUTCOME_ERROR, err
		}
		r, err := e.Evaluate(data.State())
		if err != nil {
			return EVENT_OUTCOME_ERROR, err
		}
		data.Update(v, r.(float64))
	}
	return EVENT_OUTCOME_TRUE, nil
}

func (f Event) GetID() string {
	return f.ID
}

func find(a []string, i string) bool {
	for _, v := range a {
		if v == i {
			return true
		}
	}
	return false
}

func (influenced Event) IsInfluencedBy(influencer Event) (bool, error) {
	vars := []string{}
	for _, a := range influencer.Assignments {
		v, _, err := parse(a)
		if err != nil {
			return false, err
		}
		vars = append(vars, v)
	}
	if len(vars) == 0 {
		return false, nil
	}
	for _, s := range influenced.Statements {
		e, err := govaluate.NewEvaluableExpression(s)
		if err != nil {
			return false, err
		}
		for _, v := range e.Vars() {
			if find(vars, v) {
				return true, nil
			}
		}
	}
	return false, nil
}

func (f Event) IsValid() error {
	for _, s := range f.Statements {
		_, err := govaluate.NewEvaluableExpression(s)
		if err != nil {
			return err
		}
	}
	return nil
}