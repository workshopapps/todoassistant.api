package eventEntity

type Payload struct {
	Action    string            `json:"action"`
	SubAction string            `json:"sub_action"`
	Data      map[string]string `json:"data"`
}
