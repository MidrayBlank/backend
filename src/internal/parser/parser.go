package parser

type ParserManager struct {
	Demography *DemographyParser
	Rosstat    *RosstatParser
	AgeSex     *AgeSexParser
	Migration  *MigrationParser
}

func NewParserManager() *ParserManager {
	return &ParserManager{
		Demography: NewDemographyParser(),
		Rosstat:    NewRosstatParser(),
		AgeSex:     NewAgeSexParser(),
		Migration:  NewMigrationParser(),
	}
}
