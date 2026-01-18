package utils

import "gorm.io/gorm"

// 通过结构体实现子包
type SQL struct{}

var Sql = &SQL{}

func (*SQL) Paginate(pageNo, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// Validate and normalize page number
		if pageNo <= 0 {
			pageNo = 1
		}
		// Validate and normalize page size (max 100, min 10)
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		// Calculate offset: (page - 1) * size
		// Example: page 1, size 10 -> offset 0
		//          page 2, size 10 -> offset 10
		offset := (pageNo - 1) * pageSize
		// Offset: Skip N records
		// Limit: Return at most N records
		return db.Offset(offset).Limit(pageSize)
	}
}

func (*SQL) OrderCreateAt() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at desc")
	}
}
