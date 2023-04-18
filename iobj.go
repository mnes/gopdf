package gopdf

import (
	"bytes"
)

//IObj inteface for all pdf object
type IObj interface {
	init(func() *GoPdf)
	getType() string
	getObjBuff() *bytes.Buffer
	//สร้าง ข้อมูลใน pdf
	build(objID int) error
}

type emptyObj struct {
}

func (obj emptyObj) init(func() *GoPdf) {
}
func (obj emptyObj) getType() string {
	return ""
}
func (obj emptyObj) getObjBuff() *bytes.Buffer {
	return nil
}
func (obj emptyObj) build(objID int) error {
	return nil
}
