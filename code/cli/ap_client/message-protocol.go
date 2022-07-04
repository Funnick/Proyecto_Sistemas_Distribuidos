package ap_client

type ResponseMessage struct {
	Message string `json:"message"`
}

type CreateAgentMessage struct {
	Name          string `json:"name"`
	IP            string `json:"ip"`
	Port          string `json:"port"`
	Password      string `json:"password"`
	Description   string `json:"description"`
	Documentation string `json:"documentation"`
}

type SearchAgentNameMessageResponse struct {
	AgentFound Agent  `json:"agent"`
	Message    string `json:"message"`
}

type DeleteAgentMessage struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type GetAllResponse struct {
	ResponsesFound []string `json:"responses"`
	Message        string   `json:"message"`
}

type SearchAgentNameMessage struct {
	Name string `json:"name"`
}
type SearchAgentDescMessage struct {
	Description string `json:"description"`
}

type SearchAgentMessageResponse struct {
	AgentsFound []Agent `json:"agent"`
	Message     string  `json:"message"`
}

type UpdateAgentMessage struct {
	Name             string `json:"name"`
	Password         string `json:"password"`
	NewIP            string `json:"newIp"`
	NewPort          string `json:"newPort"`
	NewPassword      string `json:"newPassword"`
	NewDescription   string `json:"newDescription"`
	NewDocumentation string `json:"newDocumentation"`
}
