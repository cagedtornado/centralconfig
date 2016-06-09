package datastores

//	The Unknown database information
type UnknownDB struct{}

func (store UnknownDB) InitStore(overwrite bool) error {
	return nil
}

func (store UnknownDB) Get(configItem ConfigItem) (ConfigItem, error) {
	return ConfigItem{}, nil
}

func (store UnknownDB) GetAllForApplication(application string) ([]ConfigItem, error) {
	return nil, nil
}

func (store UnknownDB) GetAll() ([]ConfigItem, error) {
	return nil, nil
}

func (store UnknownDB) GetAllApplications() ([]string, error) {
	return nil, nil
}

func (store UnknownDB) Set(configItem ConfigItem) (ConfigItem, error) {
	return ConfigItem{}, nil
}

func (store UnknownDB) Remove(configItem ConfigItem) error {
	return nil
}
