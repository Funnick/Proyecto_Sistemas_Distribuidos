package server

type DataBasePlatform interface {
	GetByName([]byte) ([]byte, error)
	GetByFun(string) ([]string, error)
	GetAll() ([]string, error)
	Set([]byte, string) error
	Update([]byte) error
	Delete([]byte, []byte) error
}
