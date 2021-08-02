package data

const (
	TableModtimeName      = `modtime`
	ColModtimeID          = `id`
	ColModtimeUnix        = `unix`
	ColModtimeETag        = `etag`
	ColModtimeDescription = `description`
)

type Modtime struct {
	ID          int32  `xorm:"pk 'id'"`
	Unix        int64  `xorm:"'unix'"`
	ETag        string `xorm:"'etag'"`
	Description string `xorm:"'description'"`
}

func (Modtime) TableName() string {
	return TableModtimeName
}
