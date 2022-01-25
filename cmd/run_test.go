package cmd

import (
	"bytes"
	"github.com/o-kos/dmd-cli.go/pkg/dmdintf"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRunCmd() (*cobra.Command, *bytes.Buffer, error) {
	cmd := newRunCmd()
	b := bytes.NewBufferString("")
	err := dmdintf.InitTestSystem()
	return cmd, b, err
}

//func TestRunCommand(t *testing.T) {
//	cmd, b, err := createRunCmd()
//	require.NoError(t, err)
//
//	cmd.SetArgs([]string{"Dummy"})
//	cmd.SetOut(b)
//	err = cmd.Execute()
//	require.NoError(t, err)
//
//	out, err := ioutil.ReadAll(b)
//	require.NoError(t, err)
//
//	assert.Equal(t, "COQUELET8\nDummy\nMT63\nOlivia\nPSK31\n", string(out))
//}

func TestRunCommandInvalidArgs(t *testing.T) {
	cmd, b, err := createRunCmd()
	require.NoError(t, err)

	cmd.SetArgs([]string{"Dummy"})
	cmd.SetOut(b)
	err = cmd.Execute()
	assert.EqualError(t, err, "input signal missing")

	cmd.SetArgs([]string{"Dummy", "invalid_file_name", "extra_argument"})
	cmd.SetOut(b)
	err = cmd.Execute()
	assert.EqualError(t, err, "too many arguments")

	//cmd.SetArgs([]string{"Dummy", "invalid_file_name"})
	//cmd.SetOut(b)
	//err = cmd.Execute()
	//require.EqualError(t, err, "input signal missing")

	//out, err := ioutil.ReadAll(b)
	//require.NoError(t, err)
	//assert.Equal(t, "COQUELET8\nDummy\nMT63\nOlivia\nPSK31\n", string(out))
}
