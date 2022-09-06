package types

var Chains = newChainRegistry()

type Chain struct {
	Id   int
	Name string
}

func newChainRegistry() *chainRegistry {
	unknown := Chain{Id: 0, Name: "Unknown"}
	viteMain := Chain{Id: 1, Name: "vite_main"}
	viteBuidl := Chain{Id: 9, Name: "vite_buidl"}

	return &chainRegistry{
		Unknown:   unknown,
		ViteMain:  viteMain,
		ViteBuidl: viteBuidl,
		types:     []Chain{unknown, viteMain, viteBuidl},
	}
}

type chainRegistry struct {
	Unknown   Chain
	ViteMain  Chain
	ViteBuidl Chain
	types     []Chain
}

func (r *chainRegistry) List() []Chain {
	return r.types
}

func (r *chainRegistry) GetById(id int) Chain {
	for _, v := range r.List() {
		if v.Id == id {
			return v
		}
	}
	return r.Unknown
}
