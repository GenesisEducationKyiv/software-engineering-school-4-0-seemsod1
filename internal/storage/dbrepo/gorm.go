package dbrepo

import "gorm.io/gorm"

func (m *gormDBRepo) Connection() *gorm.DB {
	return m.DB
}
