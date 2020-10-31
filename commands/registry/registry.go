package registry

import (
	"errors"

	"gobunny/commands"
)

type (
	registry struct {
		registered map[string]commands.Command
	}

	Registry interface {
		Get(string) (commands.Command, bool)

		Register(commands.Command) error
		RegisterAll(...commands.Command) error
	}
)

var (
	// ErrNameAlreadyRegistered indicates a Command's Name() is already in the registry
	ErrNameAlreadyRegistered = errors.New("command name has already been registered")

	// ErrPrefixAlreadyRegistered indicates one of a Command's Prefixes()
	// values is already in the registrt
	ErrPrefixAlreadyRegistered = errors.New("command prefix has already been registered")
)

func New() *registry {
	return &registry{
		registered: map[string]commands.Command{},
	}
}

func (r *registry) Get(name string) (commands.Command, bool) {
	command, found := r.registered[name]

	return command, found
}

func (r *registry) Register(command commands.Command) error {
	if _, found := r.registered[command.Name()]; found {
		return ErrNameAlreadyRegistered
	}

	for _, prefix := range command.Prefixes() {
		if _, found := r.registered[prefix]; found {
			return ErrPrefixAlreadyRegistered
		}
	}

	r.registered[command.Name()] = command
	for _, prefix := range command.Prefixes() {
		r.registered[prefix] = command
	}

	return nil
}

func (r *registry) RegisterAll(commands ...commands.Command) error {
	for _, command := range commands {
		if err := r.Register(command); err != nil {
			return err
		}
	}

	return nil
}
