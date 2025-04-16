package storage

type service struct {
}

func New() (*service, error) {
	return &service{}, nil
}
