package parser

// DemographyRawRecord - сырые данные из демографического файла
type DemographyRawRecord struct {
	Oktmo             string
	NotZato           int
	Region            string
	MunType           string
	Municipality      string
	Year              int
	Population        int
	AveragePopulation float64
	Deaths            int
	Births            int
	Migration         int
	MortalityRate     float64
	BirthRate         float64
	MigrationRate     float64
}

// DemographyRecord - итоговая запись для БД
type DemographyRecord struct {
	ID                  int     `db:"id"`
	Code                string  `db:"code"`                 // ОКТМО
	Year                int     `db:"year"`                 // Год
	PopulationAmount    int     `db:"population_amount"`    // Численность населения
	BirthAmount         int     `db:"birth_amount"`         // Рождения
	DeathAmount         int     `db:"death_amount"`         // Смерти
	ArrivalAmount       int     `db:"arrival_amount"`       // Прибывшие
	DepartureAmount     int     `db:"departure_amount"`     // Убывшие
	MaleAmount          int     `db:"male_amount"`          // Мужчины
	FemaleAmount        int     `db:"female_amount"`        // Женщины
	LandArea            float64 `db:"land_area"`            // Площадь
	AvgSalary           float64 `db:"avg_salary"`           // Средняя зарплата
	MedicalFacilities   int     `db:"medical_facilities"`   // Медучреждения
	SchoolsCount        int     `db:"schools_count"`        // Школы
	HousingCommissioned float64 `db:"housing_commissioned"` // Жильё
}

// RosstatRawRecord - сырые данные из файлов Росстата
type RosstatRawRecord struct {
	Oktmo string
	Year  int
	Value float64
}

// AgeSexRawRecord - сырые данные половозрастной структуры
type AgeSexRawRecord struct {
	Oktmo string
	Year  int
	Sex   string
	Value int
}

// MigrationRecord - запись миграции
type MigrationRecord struct {
	Oktmo string
	Year  int
	Value int
}
