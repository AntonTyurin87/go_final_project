package tasks

import (
	"fmt"
)

// DateConvert - пересобирает дату в удобный формат
func DateConvert(search string) (string, error) {

	result := fmt.Sprint(search[6:] + search[3:5] + search[:2])

	_, err := DateValidation(result)
	if err != nil {
		fmt.Println("Преобразовать дату не удалось ", err)
		return result, err
	}
	return result, nil
}
