package server

type DataBasePlatform interface {
	GetByName([]byte) (string, error)
	GetByFun(string) ([]string, error)
	GetAll() ([]string, error)
	Set([]byte, string) error
	Update([]byte, string) error
	Delete([]byte) error
}
