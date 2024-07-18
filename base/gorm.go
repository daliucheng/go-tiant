package base

import (
	"gorm.io/gorm"
	"time"
)

type CrudModel struct {
	CreatedAt time.Time      `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"comment:最后更新时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}

type NormalPage struct {
	No      int    // 当前第几页
	Size    int    // 每页大小
	OrderBy string `json:"orderBy"` // 排序规则
}

type Option struct {
	IsNeedCnt  bool `json:"isNeedCnt"`
	IsNeedPage bool `json:"isNeedPage"`
}

// 分页示例
func NormalPaginate(page *NormalPage) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		pageNo := 1
		if page.No > 0 {
			pageNo = page.No
		}

		pageSize := page.Size
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (pageNo - 1) * pageSize
		orderBy := "id asc"
		if len(page.OrderBy) > 0 {
			orderBy = page.OrderBy
		}
		return db.Order(orderBy).Offset(offset).Limit(pageSize)
	}
}
