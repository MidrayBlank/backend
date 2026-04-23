package parser

// Parser общий интерфейс для всех парсеров
type Parser[T any] interface {
	Parse(filePath string) ([]T, error)
}

// DemographyParserInterface интерфейс для парсера демографии
type DemographyParserInterface interface {
	ParseDemography(filePath string) ([]DemographyRawRecord, error)
}

// RosstatParserInterface интерфейс для парсера показателей Росстата
type RosstatParserInterface interface {
	ParseRosstat(filePath string) ([]RosstatRawRecord, error)
}

// AgeSexParserInterface интерфейс для парсера половозрастной структуры
type AgeSexParserInterface interface {
	ParseAgeSex(filePath string) ([]AgeSexRawRecord, error)
}

// MigrationParserInterface интерфейс для парсера миграции
type MigrationParserInterface interface {
	ParseMigration(filePath string) ([]MigrationRecord, error)
}

// DataCollectorInterface интерфейс для сборщика данных
type DataCollectorInterface interface {
	MergeAll() ([]DemographyRecord, error)
}
