package cli

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ListCommandTestSuite struct {
	suite.Suite
}

func TestListCommandTestSuite(t *testing.T) {
	suite.Run(t, new(ListCommandTestSuite))
}

func (s *ListCommandTestSuite) TestNewListCommand() {
	cmd := NewListCommand()
	s.NotNil(cmd)
	s.Equal("list", cmd.Use)
	s.Contains(cmd.Short, "List available templates and packages")
}

func (s *ListCommandTestSuite) TestListTemplatesCommand() {
	err := runListTemplatesCommand([]string{}, map[string]interface{}{})
	s.NoError(err)
}

func (s *ListCommandTestSuite) TestListPackagesCommand() {
	err := runListPackagesCommand([]string{}, map[string]interface{}{})
	s.NoError(err)
}
