package response

type Health struct {
	Status string `json:"status"`
}

func NewHealth() Health {
	return Health{Status: "ok"}
}
