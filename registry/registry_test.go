package registry_test

import (
	"gobunny/commands"
	"gobunny/registry"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

type (
	command struct {
		CommandName    string
		CommandAliases []string
	}

	registryTestSuite struct {
		suite.Suite

		command  *command
		registry registry.Registry
	}
)

func (c *command) Aliases() []string {
	return c.CommandAliases
}

func (c *command) Name() string {
	return c.CommandName
}

func (c *command) Handle(args commands.Arguments, response http.ResponseWriter, request *http.Request) error {
	return nil
}

func (c *command) Help() string {
	return "help"
}

func (c *command) Readme() string {
	return "readme"
}

func (t *registryTestSuite) SetupTest() {
	t.command = &command{
		CommandName:    "test",
		CommandAliases: []string{"t"},
	}
	t.registry = registry.New()
}

func (t *registryTestSuite) TestRegistryGetHappy() {
	err := t.registry.Register(t.command)
	t.NoError(err)

	actual, found := t.registry.Get(t.command.Name())
	t.True(found)
	t.NotNil(actual)
	t.Equal(actual.(*command), t.command)
}

func (t *registryTestSuite) TestRegisterGetNotFound() {
	actual, found := t.registry.Get(t.command.Name())
	t.False(found)
	t.Nil(actual)
}

func (t *registryTestSuite) TestRegistryRegisterHappy() {
	err := t.registry.Register(t.command)
	t.NoError(err)
}

func (t *registryTestSuite) TestRegistryRegisterDuplicateName() {
	other := &command{
		CommandName:    t.command.CommandName,
		CommandAliases: []string{},
	}

	err := t.registry.Register(t.command)
	t.NoError(err)
	err = t.registry.Register(other)
	t.Error(err)
	t.EqualError(err, registry.ErrNameAlreadyRegistered)
}

func (t *registryTestSuite) TestRegistryRegisterDuplicateAliases() {
	other := &command{
		CommandName:    "unique",
		CommandAliases: t.command.CommandAliases,
	}

	err := t.registry.Register(t.command)
	t.NoError(err)
	err = t.registry.Register(other)
	t.Error(err)
	t.EqualError(err, registry.ErrAliasAlreadyRegistered)
}

func (t *registryTestSuite) TestRegistryRegisterAllHappy() {
	other := &command{
		CommandName:    "unique",
		CommandAliases: []string{"u"},
	}
	err := t.registry.RegisterAll(t.command, other)
	t.NoError(err)
}

func (t *registryTestSuite) TestRegistryRegisterAllDuplicateName() {
	other := &command{
		CommandName:    t.command.CommandName,
		CommandAliases: []string{},
	}
	err := t.registry.RegisterAll(t.command, other)
	t.Error(err)
	t.EqualError(err, registry.ErrNameAlreadyRegistered)
}

func (t *registryTestSuite) TestRegistryRegisterAllDuplicateAliases() {
	other := &command{
		CommandName:    "unique",
		CommandAliases: t.command.CommandAliases,
	}
	err := t.registry.RegisterAll(t.command, other)
	t.Error(err)
	t.EqualError(err, registry.ErrAliasAlreadyRegistered)
}

func TestRegistryTestSuite(t *testing.T) {
	suite.Run(t, &registryTestSuite{})
}
