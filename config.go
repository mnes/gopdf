package gopdf

//Config static config
type Config struct {
	//pt , mm , cm , in
	Unit       string
	PageSize   Rect
	K          float64
	Protection PDFProtectionConfig
	LogFunc    func(msg string, args ...interface{})
}

//PDFProtectionConfig config of pdf protection
type PDFProtectionConfig struct {
	UseProtection bool
	Permissions   int
	UserPass      []byte
	OwnerPass     []byte
}
