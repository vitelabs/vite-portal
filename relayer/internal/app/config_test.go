package app

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitelabs/vite-portal/relayer/internal/types"
	"github.com/vitelabs/vite-portal/shared/pkg/util/idutil"
)

func TestConfigFile(t *testing.T) {
	p := path.Join("test", idutil.NewGuid(), types.DefaultConfigFilename)
	dir := path.Dir(p)
	defer os.RemoveAll(dir)
	expected := types.NewDefaultConfig()
	expected.HostToChainMap = map[string]string{
		"buidl.vite.net": "vite_testnet",
		"node.vite.net": "vite_mainnet",
	}
	writeConfigFile(p, &expected)
	actual := types.NewDefaultConfig()
	assert.Equal(t, 0, len(actual.HostToChainMap))
	loadConfigFromFile(p, &actual)
	files, _ := os.ReadDir(dir)
	assert.Equal(t, 1, len(files))
	assert.Equal(t, 2, len(actual.HostToChainMap))
	assert.Equal(t, actual.HostToChainMap["buidl.vite.net"], "vite_testnet")
	assert.Equal(t, actual.HostToChainMap["node.vite.net"], "vite_mainnet")
}

func TestConfigFileBackup(t *testing.T) {
	p := path.Join("test", idutil.NewGuid(), types.DefaultConfigFilename)
	dir := path.Dir(p)
	defer os.RemoveAll(dir)
	config1 := types.NewDefaultConfig()
	config1.Version = "0"
	writeConfigFile(p, &config1)
	files, _ := os.ReadDir(dir)
	assert.Equal(t, 1, len(files))
	config2 := types.NewDefaultConfig()
	loadConfigFromFile(p, &config2)
	assert.NotEqual(t, config1.Version, config2.Version)
	files, _ = os.ReadDir(dir)
	assert.Equal(t, 2, len(files))
}