package session

type PayloadKind uint8

const (
	PayloadKindMessage = iota
)

// Payload wraps actual content inside to make it possible to be passed through one api.
// TODO: This concept is still on draft. It might change a lot.
type Payload struct {
	Kind    PayloadKind
	Content any
}
