package payload

type SearchParam struct {
	Query string   `query:"q"`
	Index []string `query:"index"`
}

type HTTPResponse struct {
	StatusCode int         `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

type Response struct {
	Page         int     `json:"page"`
	Results      []Movie `json:"results"`
	TotalPages   int     `json:"total_pages"`
	TotalResults int     `json:"total_results"`
}

type Movie struct {
	Title       string `json:"title"`
	Description string `json:"overview"`
	Release     string `json:"release_date"`
	Poster      string `json:"poster_path"`
	Genre       []int  `json:"genre_ids"`
}
