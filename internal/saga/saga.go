package saga

import (
	"errors"
	"fmt"
)

type Saga struct {
	steps []Step
}

type Step struct {
	Trx        func() error
	Compensate func() error
}

func New(steps []Step) *Saga {
	return &Saga{steps: steps}
}

func (s *Saga) Run() error {
	for i := range s.steps {
		err := s.steps[i].Trx()
		if err != nil {
			err = fmt.Errorf("saga failed at step %d: %w", i, err)
			for j := i; j >= 0; j-- {
				if s.steps[j].Compensate == nil {
					continue
				}
				compensateErr := s.steps[j].Compensate()
				if compensateErr != nil {
					err = errors.Join(err, fmt.Errorf("failed compensation step %d: %w", j, compensateErr))
				}
			}
			return err
		}
	}
	return nil
}
