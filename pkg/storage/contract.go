package storage

type Status string

const (
	WorkingStatus Status = "working"
	StoppedStatus Status = "stopped"
	ErrorStatus   Status = "error"
)

type StatusConnector struct {
	Status Status
}

type Connector interface {
	Connect() error
	Close() error
	Health() StatusConnector
}
