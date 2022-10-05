package types

import (
	"github.com/vitelabs/vite-portal/shared/pkg/collections"
	"github.com/vitelabs/vite-portal/shared/pkg/util/commonutil"
)

var (
	DefaultSupportedChains = []ChainConfig{
		{Id: "1", Name: "vite_main", OfficialNodeUrl: "http://127.0.0.1:23456/"},
		{Id: "9", Name: "vite_buidl", OfficialNodeUrl: "http://127.0.0.1:23456/"},
	}
)

type Chains struct {
	idToNameMap map[string]string
	db          collections.NameObjectCollectionI[ChainConfig]
}

func NewChains(cfg []ChainConfig) *Chains {
	c := &Chains{
		idToNameMap: map[string]string{},
		db:          collections.NewNameObjectCollection[ChainConfig](),
	}
	for _, v := range cfg {
		c.idToNameMap[v.Id] = v.Name
		c.db.Add(v.Name, v)
	}
	return c
}

func (c *Chains) GetById(id string) (chain ChainConfig, found bool) {
	return c.GetByName(c.idToNameMap[id])
}

func (c *Chains) GetByName(name string) (chain ChainConfig, found bool) {
	existing := c.db.Get(name)
	if commonutil.IsEmpty(existing) {
		return *new(ChainConfig), false
	}

	return existing, true
}

func (c *Chains) Count() int {
	return len(c.idToNameMap)
}

func (c *Chains) GetEnumerator() collections.EnumeratorI[ChainConfig] {
	return c.db.GetEnumerator()
}

func (c *Chains) GetAll() []ChainConfig {
	return c.db.GetEntries()
}
