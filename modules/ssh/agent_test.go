package ssh_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/ssh"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSshAgentWithKeyPair(t *testing.T) {
	t.Parallel()

	keyPair := ssh.GenerateRSAKeyPair(t, 2048)
	sshAgent := ssh.SshAgentWithKeyPair(t, keyPair)

	// Ensure that socket directory is set in environment, and it exists.
	sockFile := filepath.Join(sshAgent.SocketDir(), "ssh_auth.sock")
	assert.FileExists(t, sockFile)

	// Assert that there's 1 key in the agent.
	keys, err := sshAgent.Agent().List()
	require.NoError(t, err)
	assert.Len(t, keys, 1)

	sshAgent.Stop()

	// Is socketDir removed as expected?
	if _, err := os.Stat(sshAgent.SocketDir()); !os.IsNotExist(err) {
		assert.FailNow(t, "ssh agent failed to remove socketDir on Stop()")
	}
}

func TestSshAgentWithKeyPairs(t *testing.T) {
	t.Parallel()

	keyPair := ssh.GenerateRSAKeyPair(t, 2048)
	keyPair2 := ssh.GenerateRSAKeyPair(t, 2048)
	sshAgent := ssh.SshAgentWithKeyPairs(t, []*ssh.KeyPair{keyPair, keyPair2})

	defer sshAgent.Stop()

	keys, err := sshAgent.Agent().List()
	require.NoError(t, err)
	assert.Len(t, keys, 2)
}

func TestMultipleSshAgents(t *testing.T) {
	t.Parallel()

	keyPair := ssh.GenerateRSAKeyPair(t, 2048)
	keyPair2 := ssh.GenerateRSAKeyPair(t, 2048)

	// Start a couple of agents.
	sshAgent := ssh.SshAgentWithKeyPair(t, keyPair)
	sshAgent2 := ssh.SshAgentWithKeyPair(t, keyPair2)

	defer sshAgent.Stop()
	defer sshAgent2.Stop()

	// Collect public keys from the agents.
	keys, err := sshAgent.Agent().List()
	require.NoError(t, err)

	keys2, err := sshAgent2.Agent().List()
	require.NoError(t, err)

	// Check that all keys match up to expected.
	assert.NotEqual(t, keys, keys2)
	assert.Equal(t, strings.TrimSpace(keyPair.PublicKey), keys[0].String())
	assert.Equal(t, strings.TrimSpace(keyPair2.PublicKey), keys2[0].String())
}
