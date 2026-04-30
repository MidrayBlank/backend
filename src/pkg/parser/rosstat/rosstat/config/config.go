package config

const (
	BURYATIA_CODE = 81
)

type Config struct {
	populationIndicator         int
	populationIndicatorBuryatia int

	SubjectCodes []int

	DownloadHTMLMaxAttempts      int
	DownloadHTMLTimeoutSeconds   int
	DownloadHTMLBatchSize        int
	DownloadHTMLTimeSleepSeconds int

	DownloadCSVMaxAttempts      int
	DownloadCSVTimeoutSeconds   int
	DownloadCSVBatchSize        int
	DownloadCSVTimeSleepSeconds int
}

func NewConfig() *Config {
	return &Config{
		populationIndicator:         8112027,
		populationIndicatorBuryatia: 8312027,
		SubjectCodes: []int{
			1, 3, 4, 5, 7, 8, 10, 11, 12, 14, 15, 17, 18, 19, 20, 22, 24,
			// 25, 26, 27, 28, 29, 30, 32, 33, 34, 35, 36, 37, 38, 40, 41, 42,
			// 44, 45, 46, 47, 49, 50, 52, 53, 54, 56, 57, 58, 60, 61, 63, 64,
			// 65, 66, 67, 68, 69, 70, 71, 73, 75, 76, 77, 78, 79, 80, 81, 82,
			// 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99,
		},
		DownloadHTMLMaxAttempts:      6,
		DownloadHTMLTimeoutSeconds:   5,
		DownloadHTMLBatchSize:        5,
		DownloadHTMLTimeSleepSeconds: 2,
		DownloadCSVMaxAttempts:       6,
		DownloadCSVTimeoutSeconds:    5,
		DownloadCSVBatchSize:         5,
		DownloadCSVTimeSleepSeconds:  3,
	}
}

func (c *Config) GetPopulationIndicator(code int) int {
	switch code {
	case BURYATIA_CODE:
		return c.populationIndicatorBuryatia
	default:
		return c.populationIndicator
	}
}
