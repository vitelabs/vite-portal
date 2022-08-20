package interfaces

type ClientI interface {
	Connect() error
	Subscribe(c chan string)
}

type OrchestratorI interface {
	GetStatus()
}