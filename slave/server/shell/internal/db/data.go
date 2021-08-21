package db

const (
	tableName = `data_of_slave_shell`

	colID         = `id`
	colUserName   = `username`
	colName       = `name`
	colFontSize   = `fontSize`
	colFontFamily = `fontFamily`
)

type DataOfSlaveShell struct {
	ID         int64  `xorm:"pk 'id'"`
	UserName   string `xorm:"index 'username' default('') "`
	Name       string `xorm:"'name' default('') "`
	FontSize   int32  `xorm:"'fontSize' default('') "`
	FontFamily string `xorm:"'fontFamily' default('') "`
}

func (DataOfSlaveShell) TableName() string {
	return tableName
}
