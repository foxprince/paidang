package repository

import "gorm.io/gorm"

// Store 数据访问层，handler 不直接碰 db
type Store struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Store { return &Store{db: db} }

func (s *Store) DB() *gorm.DB { return s.db }
