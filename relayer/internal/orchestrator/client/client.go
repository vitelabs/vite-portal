package client

type OrchestratorClient struct {
}

func NewOrchestratorClient() OrchestratorClient {
	return OrchestratorClient{}
}

func (o OrchestratorClient) Connect() error {
	return nil
}

func (o OrchestratorClient) Subscribe(c chan string) {
}
