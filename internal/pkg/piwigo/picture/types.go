package picture

type uploadChunkResponse struct {
	Status string      `json:"stat"`
	Result interface{} `json:"result"`
}

type imageExistResponse struct {
	Stat   string            `json:"stat"`
	Result map[string]string `json:"result"`
}
