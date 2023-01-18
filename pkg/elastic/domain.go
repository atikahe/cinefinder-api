package elastic

type ElasticData struct {
	Title       string `mapstructure:"title" json:"title"`
	Description string `mapstructure:"meta_description" json:"description"`
	URL         string `mapstructure:"additional_urls" json:"url"`
	Dir         string `mapstructure:"url_path_dir1" json:"dir"`
	Subdir      string `mapstructure:"url_path_dir3,omitempty" json:"subdir,omitempty"`
	// Image       string `json:"-"`
}

type ElasticHit struct {
	ID      string        `mapstructure:"_id"`
	Ignored []interface{} `mapstructure:"_ignored"`
	Index   string        `mapstructure:"_index"`
	Score   float64       `mapstructure:"_score"`
	Source  ElasticData   `mapstructure:"_source"`
}

type ElasticTotalResponse struct {
	Value float64 `mapstructure:"value"`
}

type ElasticHitParent struct {
	Hits  []ElasticHit         `mapstructure:"hits"`
	Total ElasticTotalResponse `mapstructure:"total"`
}

type ElasticResponse struct {
	Hits ElasticHitParent `mapstructure:"hits"`
	Took float64          `mapstructure:"took"`
}
