package processor

// Response ...
type Response struct {
	Key         string `json:"key"`
	ContentType string `json:"contentType"`
	Data        string `json:"data"`
	Error       string `json:"error"`
}
