package gopdf

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
)

//ImageReaderObj image object
type ImageReaderObj struct {
	buffer bytes.Buffer
	//imagepath string

	raw            []byte
	imginfo        imgInfo
	pdfProtection  *PDFProtection
	smarkObjID     int
	deviceRGBObjID int

	getRoot   func() *GoPdf
	getReader func() (io.Reader, error)
}

func (i *ImageReaderObj) init(funcGetRoot func() *GoPdf) {
	i.getRoot = funcGetRoot
}

func (i *ImageReaderObj) setProtection(p *PDFProtection) {
	i.pdfProtection = p
}

func (i *ImageReaderObj) setImageReader(funcGetReader func() (io.Reader, error)) {
	i.getReader = funcGetReader
}

func (i *ImageReaderObj) protection() *PDFProtection {
	return i.pdfProtection
}

func (i *ImageReaderObj) build(objID int) error {
	reader, err := i.getReader()
	if err != nil {
		return err
	}
	if i.raw, err = ioutil.ReadAll(reader); err != nil {
		return err
	}
	if err = i.parse(); err != nil {
		return err
	}

	if i.haveSMask() {
		i.imginfo.smarkObjID = i.smarkObjID
	}

	if i.isColspaceIndexed() {
		i.imginfo.deviceRGBObjID = i.deviceRGBObjID
	}

	buff, err := buildImgProp(i.imginfo)
	if err != nil {
		return err
	}
	_, err = buff.WriteTo(&i.buffer)
	if err != nil {
		return err
	}

	i.buffer.WriteString(fmt.Sprintf("/Length %d\n>>\n", len(i.imginfo.data))) // /Length 62303>>\n
	i.buffer.WriteString("stream\n")
	if i.protection() != nil {
		tmp, err := rc4Cip(i.protection().objectkey(objID), i.imginfo.data)
		if err != nil {
			return err
		}
		i.buffer.Write(tmp)
		i.buffer.WriteString("\n")
	} else {
		i.buffer.Write(i.imginfo.data)
	}
	i.buffer.WriteString("\nendstream\n")

	return nil
}

func (i *ImageReaderObj) isColspaceIndexed() bool {
	return isColspaceIndexed(i.imginfo)
}

func (i *ImageReaderObj) haveSMask() bool {
	return haveSMask(i.imginfo)
}

func (i *ImageReaderObj) createSMask() (*SMask, error) {
	var smk SMask
	smk.setProtection(i.protection())
	smk.w = i.imginfo.w
	smk.h = i.imginfo.h
	smk.colspace = "DeviceGray"
	smk.bitsPerComponent = "8"
	smk.filter = i.imginfo.filter
	smk.data = i.imginfo.smask
	smk.decodeParms = fmt.Sprintf("/Predictor 15 /Colors 1 /BitsPerComponent 8 /Columns %d", i.imginfo.w)
	return &smk, nil
}

func (i *ImageReaderObj) createDeviceRGB() (*DeviceRGBObj, error) {
	var dRGB DeviceRGBObj
	dRGB.data = i.imginfo.pal
	return &dRGB, nil
}

func (i *ImageReaderObj) getType() string {
	return "Image"
}

func (i *ImageReaderObj) getObjBuff() *bytes.Buffer {
	return &(i.buffer)
}

//GetRect get rect of img
func (i *ImageReaderObj) GetRect() *Rect {

	m, _, err := image.Decode(bytes.NewBuffer(i.raw))
	if err != nil {
		return nil
	}

	imageRect := m.Bounds()
	k := 1
	w := -128 //init
	h := -128 //init
	if w < 0 {
		w = -imageRect.Dx() * 72 / w / k
	}
	if h < 0 {
		h = -imageRect.Dy() * 72 / h / k
	}
	if w == 0 {
		w = h * imageRect.Dx() / imageRect.Dy()
	}
	if h == 0 {
		h = w * imageRect.Dy() / imageRect.Dx()
	}

	var rect = new(Rect)
	rect.H = float64(h)
	rect.W = float64(w)

	return rect
}

func (i *ImageReaderObj) parse() error {

	imginfo, err := parseImg(i.raw)
	if err != nil {
		return err
	}
	i.imginfo = imginfo

	return nil
}
