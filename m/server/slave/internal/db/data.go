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
	colParent      = `parent`
)

type DataOfSlave struct {
	ID          int64  `xorm:"pk autoincr 'id'"`
	Name        string `xorm:"unique 'name' default('') "`
	Description string `xorm:"'description' default('') "`
	Code        string `xorm:"unique 'code' default('') "`
	Parent      int64  `xorm:"index 'parent' default(0) "`
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
		Parent:      d.Parent,
	}
}
