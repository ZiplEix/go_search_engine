package db

import "gorm.io/gorm"

type SearchSettings struct {
	gorm.Model
	SearchOn bool `json:"searchOn"`
	AddNew   bool `json:"addNew"`
	Amount   uint `json:"amount"`
}

func (s *SearchSettings) Get() error {
	err := DBConn.Where("id = ?", 1).First(s).Error
	return err
}

func (s *SearchSettings) Update() error {
	err := DBConn.Save(s).Error
	return err
}
