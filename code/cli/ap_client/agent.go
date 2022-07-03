package ap_client

type Address struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

type Agent struct {
	Name          string  `json:"name"`
	EndPoint      Address `json:"endpoint"`
	Password      string  `json:"password"`
	Description   string  `json:"description"`
	Documentation string  `json:"documentation"`
}

func (agent Agent) Print() string {
	var print string
	name := "Name: " + agent.Name + "\n"
	endpoint := "Endpoint: " + agent.EndPoint.IP + ":" + agent.EndPoint.Port + "\n"
	desc := "Description: " + agent.Description + "\n"
	docu := "Documentación: " + agent.Documentation
	print = name + endpoint + desc + docu
	return print

}
