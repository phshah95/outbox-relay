package outbox

type OutboxEvent struct {
	Message     []byte
	Id          string
	Destination string
}
