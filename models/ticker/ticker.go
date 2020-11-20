package ticker

import (
	"encoding/json"
	"fmt"
	"gobunny/store/key"
	"gobunny/store/model"
)

// Ticker maps a particular securities ticker to an appropriate MarketWatch page
type Ticker struct {
	Symbol string `json:"symbol"`
	Href   string `json:"href"`
}

// New returns a new Ticker instance with the given MutatorFns applied
func New(mutators ...MutatorFn) (*Ticker, error) {
	t := &Ticker{}

	for _, mutator := range mutators {
		if err := mutator(t); err != nil {
			return nil, err
		}
	}

	return t, nil
}

// Marshal implements store.Marshaler
func (t *Ticker) Marshal() (string, error) {
	buf, err := json.Marshal(t)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

// Unmarshal implements store.Unmarshaler
func (t *Ticker) Unmarshal(buf string) error {
	if err := json.Unmarshal([]byte(buf), t); err != nil {
		return err
	}

	if err := t.Validate(); err != nil {
		return err
	}

	return nil
}

// Key implements store.Model
func (t *Ticker) Key() model.Key {
	return key.Must(Key(t.Symbol))
}

// Value implements store.Model
func (t *Ticker) Value() model.Value {
	value, err := t.Marshal()
	if err != nil {
		panic(fmt.Errorf("corrupt data structure: %s", err.Error()))
	}

	return model.Value(value)
}

// Validate implements store.Model
func (t *Ticker) Validate() error {
	mutators := []MutatorFn{
		WithSymbol(t.Symbol),
		WithHref(t.Href),
	}

	for _, mutator := range mutators {
		if err := mutator(t); err != nil {
			return err
		}
	}

	return nil
}
