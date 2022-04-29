package socket

type Output struct {
	BaseData

	uuid    string
	Method  string `json:"method"`
	Channel string `json:"channel"`
	Code    int    `json:"code"`
}
