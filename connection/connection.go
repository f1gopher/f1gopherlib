package connection

type Payload struct {
	Name      string
	Data      []byte
	Timestamp string
}

type Connection interface {
	Connect() (error, <-chan Payload)
}
