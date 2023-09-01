package response

type Response struct {
	MessageID int64  `json:"messageid"`
	Message   string `json:"message"`
	Status    int    `json:"status"`
}
