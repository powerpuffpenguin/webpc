package db

import (
	grpc_slave "github.com/powerpuffpenguin/webpc/protocol/slave"
)

const (
	tableName = `data_of_slave`

	colID          = `id`
	colName        = `name`
	colDescription = `description`
	colCode        = `code`
)

type DataOfSlave struct {
	ID          int64  `xorm:"pk autoincr 'id'"`
	Name        string `xorm:"unique 'name'"`
	Description string `xorm:"'description'"`
	Code        string `xorm:"unique 'code'"`
}

func (DataOfSlave) TableName() string {
	return tableName
}

func (d *DataOfSlave) ToPB() *grpc_slave.Data {
	return &grpc_slave.Data{
		Id:          d.ID,
		Name:        d.Name,
		Description: d.Description,
		Code:        d.Code,
	}
}
