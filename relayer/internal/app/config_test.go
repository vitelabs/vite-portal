package app

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitelabs/vite-portal/internal/types"
	"github.com/vitelabs/vite-portal/internal/util/idutil"
)

func TestConfigFile(t *testing.T) {
	filename := fmt.Sprintf("%s.json", idutil.NewGuid())
	expected := types.NewDefaultConfig()
	expected.HostToChainMap = map[string]string{
		"buidl.vite.net": "vite_testnet",
		"node.vite.net": "vite_mainnet",
	}
	writeConfigFile(filename, &expected)
	actual := types.NewDefaultConfig()
	assert.Equal(t, 0, len(actual.HostToChainMap))
	loadConfigFromFile(filename, &actual)
	os.Remove(filename)
	assert.Equal(t, 2, len(actual.HostToChainMap))
	assert.Equal(t, expected.HostToChainMap["buidl.vite.net"], "vite_testnet")
	assert.Equal(t, expected.HostToChainMap["node.vite.net"], "vite_mainnet")
}
