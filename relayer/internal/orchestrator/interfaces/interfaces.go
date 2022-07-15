package interfaces

type ClientI interface {
	Connect() error
	Subscribe(c chan string)
}