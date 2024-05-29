package infra

type MessageBrokerClient interface {
	Close() error
}