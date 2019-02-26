package category

type PiwigoCategory struct {
	Id       int
	ParentId int
	Name     string
	Key      string
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

type createCategoryResponse struct {
	Status  string `json:"stat"`
	Err     int    `json:"err"`
	Message string `json:"message"`
	Result  struct {
		Info string `json:"info"`
		ID   int    `json:"id"`
	} `json:"result"`
}
