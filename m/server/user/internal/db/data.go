package db

import (
	grpc_user "github.com/powerpuffpenguin/webpc/protocol/user"
)

const (
	tableName = `data_of_user`

	colID            = `id`
	colName          = `name`
	colNickname      = `nickname`
	colPassword      = `password`
	colAuthorization = `authorization`
)

type DataOfUser struct {
	ID            int64   `xorm:"pk autoincr 'id'"`
	Name          string  `xorm:"unique 'name' default('') "`
	Nickname      string  `xorm:"'nickname' default('') "`
	Password      string  `xorm:"'password' default('') "`
	Authorization []int32 `xorm:"'authorization'"`
}

func (DataOfUser) TableName() string {
	return tableName
}

func (d *DataOfUser) ToPB() *grpc_user.Data {
	return &grpc_user.Data{
		Id:            d.ID,
		Name:          d.Name,
		Nickname:      d.Nickname,
		Authorization: d.Authorization,
	}
}
