package responses

type Responses struct {
	Code        int64       `json:"code"`
	Description string      `json:"description"`
	Data        interface{} `json:"data"`
}
