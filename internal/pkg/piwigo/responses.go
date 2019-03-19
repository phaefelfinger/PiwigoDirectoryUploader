package piwigo

type responseStatuser interface {
	responseStatus() string
}

type loginResponse struct {
	Status      string `json:"stat"`
	Result      bool   `json:"result"`
	ErrorNumber int    `json:"err"`
	Message     string `json:"message"`
}

func (r loginResponse) responseStatus() string {
	return r.Status
}

type getStatusResponse struct {
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

func (r getStatusResponse) responseStatus() string {
	return r.Status
}

type logoutResponse struct {
	Status string `json:"stat"`
	Result bool   `json:"result"`
}

func (r logoutResponse) responseStatus() string {
	return r.Status
}

type getCategoryListResponse struct {
	Status string `json:"stat"`
	Result struct {
		Categories []struct {
			ID                      int    `json:"id"`
			Name                    string `json:"name"`
			Comment                 string `json:"comment,omitempty"`
			Permalink               string `json:"permalink,omitempty"`
			Status                  string `json:"status,omitempty"`
			Uppercats               string `json:"uppercats,omitempty"`
			GlobalRank              string `json:"global_rank,omitempty"`
			IDUppercat              int    `json:"id_uppercat,string,omitempty"`
			NbImages                int    `json:"nb_images,omitempty"`
			TotalNbImages           int    `json:"total_nb_images,omitempty"`
			RepresentativePictureID string `json:"representative_picture_id,omitempty"`
			DateLast                string `json:"date_last,omitempty"`
			MaxDateLast             string `json:"max_date_last,omitempty"`
			NbCategories            int    `json:"nb_categories,omitempty"`
			URL                     string `json:"url,omitempty"`
			TnURL                   string `json:"tn_url,omitempty"`
		} `json:"categories"`
	} `json:"result"`
}

func (r getCategoryListResponse) responseStatus() string {
	return r.Status
}

type createCategoryResponse struct {
	Status  string `json:"stat"`
	Err     int    `json:"err"`
	Message string `json:"message"`
	Result  struct {
		Info string `json:"info"`
		ID   int    `json:"id"`
	} `json:"result"`
}

func (r createCategoryResponse) responseStatus() string {
	return r.Status
}

type uploadChunkResponse struct {
	Status string      `json:"stat"`
	Result interface{} `json:"result"`
}

func (r uploadChunkResponse) responseStatus() string {
	return r.Status
}

type fileAddResponse struct {
	Status string `json:"stat"`
	Result struct {
		ImageID int    `json:"image_id"`
		URL     string `json:"url"`
	} `json:"result"`
}

func (r fileAddResponse) responseStatus() string {
	return r.Status
}

type imageExistResponse struct {
	Status string            `json:"stat"`
	Result map[string]string `json:"result"`
}

func (r imageExistResponse) responseStatus() string {
	return r.Status
}

type checkFilesResponse struct {
	Status string            `json:"stat"`
	Result map[string]string `json:"result"`
}

func (r checkFilesResponse) responseStatus() string {
	return r.Status
}
