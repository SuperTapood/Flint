package base

type Connection interface {
	List() map[string]any
	Deploy(map[string]any, string)
}
