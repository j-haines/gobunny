package registry

import (
	"errors"

	"gobunny/commands"
)

type (
	registry struct {
		registered map[string]commands.Command
	}

	// Registry is responsible for the registration and storing of Commands
	Registry interface {
		// Get returns a Command registered under the given alias
		Get(string) (commands.Command, bool)

		// Register adds a Command to the Registry store under its Aliases
		Register(commands.Command) error

		// RegisterAll adds a slice of Commands to the Registry store, each
		// under their corresponding Aliases
		RegisterAll(...commands.Command) error
	}
)

var (
	// ErrAliasAlreadyRegistered indicates one of a Command's Aliases()
	// values is already in the registry
	ErrAliasAlreadyRegistered = "command alias has already been registered"

	// ErrNameAlreadyRegistered indicates a Command's Name() value is already in the registry
	ErrNameAlreadyRegistered = "command name has already been registered"
)

// New returns an instance of a Registry implementation
func New() Registry {
	return &registry{
		registered: map[string]commands.Command{},
	}
}

// Get implements Registry
func (r *registry) Get(name string) (commands.Command, bool) {
	command, found := r.registered[name]

	return command, found
}

// Register implements Registry
func (r *registry) Register(command commands.Command) error {
	if _, found := r.registered[command.Name()]; found {
		return errors.New(ErrNameAlreadyRegistered)
	}

	for _, alias := range command.Aliases() {
		if _, found := r.registered[alias]; found {
			return errors.New(ErrAliasAlreadyRegistered)
		}
	}

	r.registered[command.Name()] = command
	for _, alias := range command.Aliases() {
		r.registered[alias] = command
	}

	return nil
}

// RegisterAll implements Registry
func (r *registry) RegisterAll(commands ...commands.Command) error {
	for _, command := range commands {
		if err := r.Register(command); err != nil {
			return err
		}
	}

	return nil
}
