package picture

type uploadChunkResponse struct {
	Status string      `json:"stat"`
	Result interface{} `json:"result"`
}

type fileAddResponse struct {
	Status string `json:"stat"`
	Result struct {
		ImageID int    `json:"image_id"`
		URL     string `json:"url"`
	} `json:"result"`
}

type imageExistResponse struct {
	Stat   string            `json:"stat"`
	Result map[string]string `json:"result"`
}
