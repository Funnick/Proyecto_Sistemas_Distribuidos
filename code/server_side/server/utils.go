package server

type DBChord interface {
	Delete(string, string) error
	GetByName(string) (string, error)
	GetByFunction(string) ([]string, error)
	Set(Agent) error
	Update(string, string, string) error
}
