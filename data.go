package berrors

type BadInput struct {
	Fields []field `json:"fields"`
}

type field struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

type KafkaErrorMessage struct {
	Code    int    `json:"code"`
	Content []byte `json:"content"`
}
