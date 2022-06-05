package server

type ResponseMessage struct {
	Message string `json:"message"`
}

type CreateAgentMessage struct {
	IP            string `json:"ip"`
	Port          string `json:"port"`
	Password      string `json:"password"`
	Description   string `json:"description"`
	Documentation string `json:"documentation"`
}

type DeleteAgentMessage struct {
	AID      string `json:"aid"`
	Password string `json:"password"`
}

type SearchAgentMessage struct {
	Description string `json:"description"`
}

type SearchAgentMessageResponse struct {
	AgentFound Agent  `json:"agent"`
	Message    string `json:"message"`
}

type UpdateAgentMessage struct {
	AID              string `json:"aid"`
	Password         string `json:"password"`
	NewIP            string `json:"newIp"`
	NewPort          string `json:"newPort"`
	NewPassword      string `json:"newPassword"`
	NewDescription   string `json:"newDescription"`
	NewDocumentation string `json:"newDocumentation"`
}
