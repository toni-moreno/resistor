package config

// SQLConfig read from DB
type SQLConfig struct {
	dummy map[string]int
}

func Init(cfg *SQLConfig) error {

	log.Debug("--------------------Initializing Config-------------------")

	log.Debug("-----------------------END Config metrics----------------------")
	return nil
}
