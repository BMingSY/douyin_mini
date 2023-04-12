package mysql

import "tiktok/initial"

type Test struct {
	Id   int
	Name string
}

func (Test) TableName() string {
	return "test"
}

func GetNameByID(id int) (string, error) {
	db := initial.Database
	var book Test
	err := db.Where("id=?", id).First(&book).Error
	if err != nil {
		return "", err
	}
	return book.Name, nil

}
