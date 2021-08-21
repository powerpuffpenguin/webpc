package db

import (
	"github.com/powerpuffpenguin/webpc/db/manipulator"
)

func Add(data *DataOfSlaveShell) (e error) {
	_, e = manipulator.Engine().InsertOne(data)
	if e != nil {
		return
	}
	return
}
func Remove(id int64) (changed bool, e error) {
	rows, e := manipulator.Engine().ID(id).Delete(&DataOfSlaveShell{})
	changed = rows != 0
	return
}
func Rename(id int64, name string) (changed bool, e error) {
	rows, e := manipulator.Engine().ID(id).Cols(colName).Update(&DataOfSlaveShell{
		Name: name,
	})
	changed = rows != 0
	return
}
func FontSize(id int64, fontSize int32) (changed bool, e error) {
	rows, e := manipulator.Engine().ID(id).Cols(colFontSize).Update(&DataOfSlaveShell{
		FontSize: fontSize,
	})
	changed = rows != 0
	return
}
func FontFamily(id int64, fontFamily string) (changed bool, e error) {
	rows, e := manipulator.Engine().ID(id).Cols(colFontFamily).Update(&DataOfSlaveShell{
		FontFamily: fontFamily,
	})
	changed = rows != 0
	return
}
func List() (items []DataOfSlaveShell, e error) {
	e = manipulator.Engine().Find(&items)
	return
}
