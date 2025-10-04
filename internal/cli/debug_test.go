package cli

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type DebugCommandTestSuite struct {
	suite.Suite
}

func TestDebugCommandTestSuite(t *testing.T) {
	suite.Run(t, new(DebugCommandTestSuite))
}

func (s *DebugCommandTestSuite) TestNewDebugCommand() {
	cmd := NewDebugCommand()
	s.NotNil(cmd)
	s.Equal("debug", cmd.Use)
	s.Contains(cmd.Short, "Debug Forge OS projects")
}

func (s *DebugCommandTestSuite) TestDebugCommand() {
	err := runDebugCommand([]string{}, map[string]interface{}{
		"config": true,
		"env":    true,
		"system": true,
	})
	s.NoError(err)
}
