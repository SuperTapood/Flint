package base

type Connection interface {
	List()
	Deploy(map[string]any, string)
}
