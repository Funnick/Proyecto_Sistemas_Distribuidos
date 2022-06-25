package chord

type DataBasePlatform interface {
	Get([]byte) (string, error)
	Set([]byte, string) error
	Update([]byte, string) error
	Delete([]byte) error
}
