package configutil

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
	"github.com/vitelabs/vite-portal/shared/pkg/util/idutil"
)

func TestInitConfig(t *testing.T) {
	p := path.Join("test", idutil.NewGuid(), "test.json")
	dir := path.Dir(p)
	defer os.RemoveAll(dir)
	actual := newTestConfig()
	actual.Debug = false
	err := InitConfig(&actual, true, p, actual.Version)
	assert.NoError(t, err)
	assert.Equal(t, true, actual.Debug)
}

func TestConfigFile(t *testing.T) {
	p := path.Join("test", idutil.NewGuid(), "test.json")
	dir := path.Dir(p)
	defer os.RemoveAll(dir)
	expected := newTestConfig()
	expected.HostToChainMap = map[string]string{
		"buidl.vite.net": "vite_testnet",
		"node.vite.net":  "vite_mainnet",
	}
	WriteConfigFile(p, &expected)
	actual := newTestConfig()
	assert.Equal(t, 0, len(actual.HostToChainMap))
	LoadConfigFromFile(p, actual.Version, &actual)
	files, _ := os.ReadDir(dir)
	assert.Equal(t, 1, len(files))
	assert.Equal(t, 2, len(actual.HostToChainMap))
	assert.Equal(t, actual.HostToChainMap["buidl.vite.net"], "vite_testnet")
	assert.Equal(t, actual.HostToChainMap["node.vite.net"], "vite_mainnet")
}

func TestConfigFileBackup(t *testing.T) {
	p := path.Join("test", idutil.NewGuid(), "test.json")
	dir := path.Dir(p)
	defer os.RemoveAll(dir)
	config1 := newTestConfig()
	config1.Version = "0"
	WriteConfigFile(p, &config1)
	files, _ := os.ReadDir(dir)
	assert.Equal(t, 1, len(files))
	config2 := newTestConfig()
	LoadConfigFromFile(p, "v0.2", &config2)
	assert.NotEqual(t, config1.Version, config2.Version)
	files, _ = os.ReadDir(dir)
	assert.Equal(t, 2, len(files))
}

type testConfig struct {
	Version        string            `json:"version"`
	Debug          bool              `json:"debug"`
	HostToChainMap map[string]string `json:"hostToChainMap"`
}

func (c *testConfig) GetVersion() string {
	return c.Version
}

func (c *testConfig) GetDebug() bool {
	return c.Debug
}

func (c *testConfig) SetDebug(debug bool) {
	c.Debug = debug
}

func (c *testConfig) GetLoggingConfig() sharedtypes.LoggingConfig {
	return *new(sharedtypes.LoggingConfig)
}

func (c *testConfig) Validate() error {
	return nil
}

func newTestConfig() testConfig {
	return testConfig{
		Version:        "v0.1",
		Debug:          true,
		HostToChainMap: map[string]string{},
	}
}
