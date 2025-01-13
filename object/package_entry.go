package object

type PackageEntry struct {
	Hash     []byte `json:"hash,omitempty"`
	Content  []byte `json:"content,omitempty"`
	Ext      string `json:"ext,omitempty"`
	Filename string `json:"filename,omitempty"`
}

func NewPackageEntry(hash []byte, content []byte, ext string, filename string) *PackageEntry {
	return &PackageEntry{
		Hash:     hash,
		Content:  content,
		Ext:      ext,
		Filename: filename,
	}
}
