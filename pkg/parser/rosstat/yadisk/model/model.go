package model

type DemographyRecord struct {
	ID                  int
	Code                string   // ОКТМО
	Year                int      // Год
	PopulationAmount    int      // Численность населения
	BirthAmount         *int     // Рождения
	DeathAmount         *int     // Смерти
	ArrivalAmount       *int     // Прибывшие
	DepartureAmount     *int     // Убывшие
	MaleAmount          *int     // Мужчины (может быть NULL, если есть данные по возрастам)
	FemaleAmount        *int     // Женщины (может быть NULL)
	LandArea            *float64 // Площадь
	AvgSalary           *float64 // Средняя зарплата
	MedicalFacilities   *int     // Медучреждения
	SchoolsCount        *int     // Школы
	HousingCommissioned *float64 // Жильё
}

// RosstatAgeRecord - запись для таблицы rosstat_age (детализация по возрастам)
type RosstatAgeRecord struct {
	RosstatID    int    // FK на основную таблицу
	Age          string // Возрастная группа (например, "0-4", "5-9", "20-24")
	MaleAmount   int
	FemaleAmount int
}
