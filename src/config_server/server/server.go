package server

type ConfigServer interface {
	Start(int, string, string) error
}
