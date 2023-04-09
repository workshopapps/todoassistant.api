package dataEntity

type Country struct {
	CountryId   int    `json:"country_id"`
	CountryName string `json:"country_name"`
	TelCode     string `json:"tel_code"`
}
