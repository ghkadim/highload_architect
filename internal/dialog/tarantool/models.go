package tarantool

type message struct {
	ID   int64  `json:"id"`
	From string `json:"from"`
	To   string `json:"to"`
	Text string `json:"text"`
}
