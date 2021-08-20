package db

import (
	"github.com/powerpuffpenguin/webpc/db/manipulator"
)

func Add(data *DataOfSingleShell) (e error) {
	_, e = manipulator.Engine().InsertOne(data)
	if e != nil {
		return
	}
	return
}
func Remove(id int64) (changed bool, e error) {
	rows, e := manipulator.Engine().ID(id).Delete(&DataOfSingleShell{})
	changed = rows != 0
	return
}
func Rename(id int64, name string) (changed bool, e error) {
	rows, e := manipulator.Engine().ID(id).Cols(colName).Update(&DataOfSingleShell{
		Name: name,
	})
	changed = rows != 0
	return
}
func FontSize(id int64, fontSize int32) (changed bool, e error) {
	rows, e := manipulator.Engine().ID(id).Cols(colFontSize).Update(&DataOfSingleShell{
		FontSize: fontSize,
	})
	changed = rows != 0
	return
}
func FontFamily(id int64, fontFamily string) (changed bool, e error) {
	rows, e := manipulator.Engine().ID(id).Cols(colFontFamily).Update(&DataOfSingleShell{
		FontFamily: fontFamily,
	})
	changed = rows != 0
	return
}
