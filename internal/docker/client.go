package docker

type client struct {
	docker *client.Client
}

func New() *client {
	return &client{}
}
