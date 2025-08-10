package cfg

const (
	EncryptedExt  = ".kpenc"
	KeyFileFormat = "kpkey.%s"
	CfgFileFormat = "kpcfg.%s"
)

type NewCfgOpts struct {
	Dir     string
	Profile string
}
