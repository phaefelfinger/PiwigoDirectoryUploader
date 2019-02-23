package authentication

import "net/http/cookiejar"

type PiwigoContext struct {
	Url      string
	Username string
	Password string
	Cookies  *cookiejar.Jar
}

type LoginResponse struct {
	Status      string `json:"stat"`
	Result      bool   `json:"result"`
	ErrorNumber int    `json:"err"`
	Message     string `json:"message"`
}

type GetStatusResponse struct {
	Status string `json:"stat"`
	Result struct {
		Username            string   `json:"username"`
		Status              string   `json:"status"`
		Theme               string   `json:"theme"`
		Language            string   `json:"language"`
		PwgToken            string   `json:"pwg_token"`
		Charset             string   `json:"charset"`
		CurrentDatetime     string   `json:"current_datetime"`
		Version             string   `json:"version"`
		AvailableSizes      []string `json:"available_sizes"`
		UploadFileTypes     string   `json:"upload_file_types"`
		UploadFormChunkSize int      `json:"upload_form_chunk_size"`
	} `json:"result"`
}

type LogoutResponse struct {
	Status string `json:"stat"`
	Result bool   `json:"result"`
}
