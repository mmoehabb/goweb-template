package users

type DataModel struct {
	Username string
	Password string
}

func parseRow(row []any) DataModel {
	return DataModel{
		Username: row[0].(string),
		Password: row[1].(string),
	}
}

func parseModel(m *DataModel) map[string]string {
	var modelMap = make(map[string]string)
	if m.Username != "" {
		modelMap["Username"] = m.Username
	}
	if m.Password != "" {
		modelMap["Username"] = m.Password
	}
	return modelMap
}
