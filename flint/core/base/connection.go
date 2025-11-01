package base

type Connection interface {
	GetCurrentRevision(string) int
	List() []Deployment
	Deploy(map[string]any)
}
