package cli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/suite"
)

type VersionCommandTestSuite struct {
	suite.Suite
}

func TestVersionCommandTestSuite(t *testing.T) {
	suite.Run(t, new(VersionCommandTestSuite))
}

func (s *VersionCommandTestSuite) TestVersionCommandCreation() {
	cmd := NewVersionCommand()
	s.NotNil(cmd)
	s.Equal("version", cmd.Use)
	s.Contains(cmd.Short, "Show version information")
}

func (s *VersionCommandTestSuite) TestVersionCommandBasic() {
	cmd := NewVersionCommand()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)

	err := cmd.Execute()
	s.NoError(err)

	output := buf.String()
	s.Contains(output, "Forge OS dev")
	s.Contains(output, "Go Version:")
	s.Contains(output, "Platform:")
}

func (s *VersionCommandTestSuite) TestVersionCommandVerbose() {
	cmd := NewVersionCommand()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)

	cmd.Flags().Set("verbose", "true")
	err := cmd.Execute()
	s.NoError(err)

	output := buf.String()
	s.Contains(output, "Forge OS dev")
	s.Contains(output, "Go Version:")
	s.Contains(output, "Platform:")
	s.Contains(output, "Buildroot Version:")
	s.Contains(output, "Kernel Version:")
}
