package model

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        int64          `gorm:"primarykey;comment: 主键ID"`
	CreatedAt time.Time      `gorm:"colum:created_at ;comment: 创建时间"`
	UpdatedAt time.Time      `gorm:"colum:updated_at ;comment: 更新时间"`
	DeletedAt gorm.DeletedAt `gorm:"comment: 删除时间"`
	IsDeleted bool           `gorm:"comment: 删除标志"`
}

//password

type User struct {
	BaseModel
	Mobile   string     `gorm:"index:idx_mobile;unique;type:varchar(11);not null;comment: 手机号"`
	Password string     `gorm:"type:varchar(200) ;not null;comment: 密码"`
	NickName string     `gorm:"type:varchar(40) ;comment: 昵称"`
	Birthday *time.Time `gorm:"type:datetime ;comment: 生日"`
	Sex      int        `gorm:"type:int ;colum:sex;default:0;comment: 性别 0未知 1男 2女"`
	Role     int        `gorm:"type:int;colum:role;default:1 ;comment: 权限 0管理员 1用户"`
}
