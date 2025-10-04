package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type DoctorCommandTestSuite struct {
	suite.Suite
	tempDir string
}

func TestDoctorCommandTestSuite(t *testing.T) {
	suite.Run(t, new(DoctorCommandTestSuite))
}

func (s *DoctorCommandTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "forge-doctor-cmd-*")
	s.Require().NoError(err)
}

func (s *DoctorCommandTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *DoctorCommandTestSuite) TestDoctorCommandCreation() {
	cmd := NewDoctorCommand()
	s.NotNil(cmd)
	s.Equal("doctor", cmd.Use)
	s.Contains(cmd.Short, "Check system and diagnose issues")
}

func (s *DoctorCommandTestSuite) TestDoctorCommandBasic() {
	cmd := NewDoctorCommand()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)

	err := cmd.Execute()
	s.NoError(err)

	output := buf.String()
	s.Contains(output, "Forge OS Doctor")
	s.Contains(output, "Go Version:")
	s.Contains(output, "Platform:")
	s.Contains(output, "System Resources:")
	s.Contains(output, "Doctor check complete!")
}

func (s *DoctorCommandTestSuite) TestDoctorCommandVerbose() {
	cmd := NewDoctorCommand()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)

	cmd.Flags().Set("verbose", "true")
	err := cmd.Execute()
	s.NoError(err)

	output := buf.String()
	s.Contains(output, "Forge OS Doctor")
	s.Contains(output, "Verbose Information:")
	s.Contains(output, "GOPATH:")
	s.Contains(output, "GOROOT:")
}

func (s *DoctorCommandTestSuite) TestDoctorCommandInProject() {
	// Create a project
	projectDir := filepath.Join(s.tempDir, "test-project")
	err := createProjectStructure(projectDir, "minimal", "x86_64")
	s.NoError(err)

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(projectDir)

	cmd := NewDoctorCommand()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)

	err = cmd.Execute()
	s.NoError(err)

	output := buf.String()
	s.Contains(output, "Forge project detected")
}

func (s *DoctorCommandTestSuite) TestDoctorCommandNotInProject() {
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(s.tempDir)

	cmd := NewDoctorCommand()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)

	err := cmd.Execute()
	s.NoError(err)

	output := buf.String()
	s.Contains(output, "Not in a Forge project directory")
}
