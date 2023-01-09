package types

var ConnectionTypes = newConnectionTypeRegistry()

func newConnectionTypeRegistry() *connectionTypeRegistry {
	unknown := "unknown"
	relayer := "relayer"
	node := "node"

	return &connectionTypeRegistry{
		Unknown: unknown,
		Relayer: relayer,
		Node:    node,
		types:   []string{unknown, relayer, node},
	}
}

type connectionTypeRegistry struct {
	Unknown string
	Relayer string
	Node    string
	types   []string
}

func (r *connectionTypeRegistry) List() []string {
	return r.types
}
