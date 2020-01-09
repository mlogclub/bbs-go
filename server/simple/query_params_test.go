package simple

import (
	"database/sql"
	"fmt"
	"testing"
)

type User struct {
	GormModel
	Username    sql.NullString `gorm:"size:32;unique;" json:"username" form:"username"`
	Email       sql.NullString `gorm:"size:128;unique;" json:"email" form:"email"`
	Nickname    string         `gorm:"size:16;" json:"nickname" form:"nickname"`
	Avatar      string         `gorm:"type:text" json:"avatar" form:"avatar"`
	Password    string         `gorm:"size:512" json:"password" form:"password"`
	Status      int            `gorm:"index:idx_status;not null" json:"status" form:"status"`
	Roles       string         `gorm:"type:text" json:"roles" form:"roles"`
	Type        int            `gorm:"not null" json:"type" form:"type"`
	Description string         `gorm:"type:text" json:"description" form:"description"`
	CreateTime  int64          `json:"createTime" form:"createTime"`
	UpdateTime  int64          `json:"updateTime" form:"updateTime"`
}

func TestQueryParams(t *testing.T) {
	if err := OpenDB("root:123456@tcp(localhost:3306)/mlog_db2?charset=utf8mb4&parseTime=True&loc=Local", 5, 20, true); err != nil {
		panic(err)
	}
	// var users []User
	// _ = NewQueryParams(nil).Desc("id").Find(DB(), &users)
	// fmt.Println(FormatJson(users))

	// count, _ := NewQueryParams(nil).Desc("id").Count(DB(), &User{})
	// fmt.Println(count)

	// params.Query(db).Find(&list)
	// params.Count(db).Model(&model.Article{}).Count(&params.Paging.Total)
	// NewSqlCnd().Where("username = ? or email = ?", "username", "email").Where("password = ?", 123).Query(DB()).Find(&users)

	// var users []User
	// NewSqlCnd().In("id", []int64{1, 2, 3}).Find(db, &users)
	//
	// fmt.Println(len(users))
	// for _, user := range users {
	// 	fmt.Println(user.Nickname)
	// }

	var users []User
	NewSqlCnd().Cols("id", "status").In("id", []int64{1, 2, 3}).Find(db, &users)

	for _, user := range users {
		fmt.Println(user.Id, user.Nickname)
	}
}
