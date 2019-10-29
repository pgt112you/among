package server

type CommonServerInfo interface {
	Unmarshal([]byte) error
}

type CommonServer interface {
	Run()
}
