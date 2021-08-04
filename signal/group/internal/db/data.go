package db

const (
	tableName = `data_of_group`

	colID          = `id`
	colParent      = `parent`
	colName        = `name`
	colDescription = `description`
)

type DataOfGroup struct {
	ID          int64  `xorm:"pk autoincr 'id'"`
	Parent      int64  `xorm:"'parent' default(0)"`
	Name        string `xorm:"unique 'name' default('') "`
	Description string `xorm:"'description' default('') "`
}

func (DataOfGroup) TableName() string {
	return tableName
}
