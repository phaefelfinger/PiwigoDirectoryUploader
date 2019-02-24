package category

type PiwigoCategory struct {
	Id   int
	Name string
	Key  string
}

type getCategoryListResponse struct {
	Status   string `json:"stat"`
	Result struct {
		Categories []struct {
			ID                      int         `json:"id"`
			Name                    string      `json:"name"`
			Comment                 string      `json:"comment"`
			Permalink               interface{} `json:"permalink"`
			Status                  string      `json:"status"`
			Uppercats               string      `json:"uppercats"`
			GlobalRank              string      `json:"global_rank"`
			IDUppercat              interface{} `json:"id_uppercat"`
			NbImages                int         `json:"nb_images"`
			TotalNbImages           int         `json:"total_nb_images"`
			RepresentativePictureID string      `json:"representative_picture_id"`
			DateLast                interface{} `json:"date_last"`
			MaxDateLast             string      `json:"max_date_last"`
			NbCategories            int         `json:"nb_categories"`
			URL                     string      `json:"url"`
			TnURL                   string      `json:"tn_url"`
		} `json:"categories"`
	} `json:"result"`
}
