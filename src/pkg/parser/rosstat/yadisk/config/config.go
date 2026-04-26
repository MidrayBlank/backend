package config

type Config struct {
	PopulationURL          string
	BirthURL               string
	DeathURL               string
	ArrivalURL             string
	DepartureURL           string
	LandAreaURL            string
	AverageSalaryURL       string
	MedicialFacilitiesURL  string
	SchoolsURL             string
	HousingCommissionedURL string
	MaleFemaleAgeURL       string
}

func NewConfig() *Config {
	return &Config{
		PopulationURL:          "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_bdmo_118_v20250918/indicators/section31/data_Y48112027_112_v20250918.zip",
		BirthURL:               "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_bdmo_118_v20250918/indicators/section31/data_Y48112003_112_v20250918.zip",
		DeathURL:               "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_bdmo_118_v20250918/indicators/section31/data_Y48112001_112_v20250918.zip",
		ArrivalURL:             "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_bdmo_118_v20250918/indicators/section31/data_Y48112021_112_v20250918.zip",
		DepartureURL:           "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_bdmo_118_v20250918/indicators/section31/data_Y48112022_112_v20250918.zip",
		LandAreaURL:            "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_bdmo_118_v20250918/indicators/section6/data_Y48006001_112_v20250918.zip",
		AverageSalaryURL:       "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_bdmo_118_v20250918/indicators/section32/data_Y48423007_112_v20250918.zip",
		MedicialFacilitiesURL:  "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_bdmo_118_v20250918/indicators/section33/data_Y48018000_112_v20250918.zip",
		SchoolsURL:             "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_bdmo_118_v20250918/indicators/section34/data_Y48015001_112_v20250918.zip",
		HousingCommissionedURL: "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_bdmo_118_v20250918/indicators/section37/data_Y48010001_112_v20250918.zip",
		MaleFemaleAgeURL:       "https://storage.yandexcloud.net/tochno-st-catalog/Rosstat/data_bdmo_118_v20250918/indicators/section31/data_Y48112014_112_v20250918.zip",
	}
}
