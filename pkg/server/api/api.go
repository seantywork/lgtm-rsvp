package api

type CLIENT_REQ struct {
	Data string `json:"data"`
}

type SERVER_RESP struct {
	Status string `json:"status"`
	Reply  string `json:"reply"`
}
