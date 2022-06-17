package server

type Address struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

type Agent struct {
	Name          string   `json:"name"`
	AID           []byte   `json:"aid"`
	EndPoint      *Address `json:"endpoint"`
	Password      string   `json:"password"`
	Description   string   `json:"description"`
	Documentation string   `json:"documentation"`
	DesSplit      []string
}
