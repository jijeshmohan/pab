package pab

import "errors"

// Adapter interface defines chat adapters
type Adapter interface {
	Run()
	Stop()
	//
	Receive(*Message)
	Send(*Response, string)
	Name() string
}

type initAdapter func(*Bot) (Adapter, error)

var adapters = map[string]initAdapter{}

func newAdapter(name string, b *Bot) (Adapter, error) {
	initAdapter, ok := adapters[name]
	if !ok {
		return nil, errors.New("Unable to find adapter with name " + name)
	}
	return initAdapter(b)
}

// RegisterAdapter is for registering a new adapter.
func RegisterAdapter(name string, f func(*Bot) (Adapter, error)) {
	adapters[name] = f
}
