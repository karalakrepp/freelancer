package helper

import "errors"

func SettingInQueryWithID(dbName string, param string, column string) (string, error) {
	arr := SplitComma(param)
	lengthArr := len(arr)
	initialQuery := "SELECT " + column + " FROM " + dbName + " WHERE id IN ("
	for i := 0; i < lengthArr; i++ {
		initialQuery = initialQuery + arr[i] + ","
	}

	var lengthString = len(initialQuery)

	if lengthArr > 0 && lengthString > 0 && initialQuery[lengthString-1] == ',' {
		initialQuery = initialQuery[:lengthString-1]
		initialQuery += ")"
	} else if lengthArr <= 0 {
		return "", errors.New("Parameter is empty")
	}

	return initialQuery, nil
}
