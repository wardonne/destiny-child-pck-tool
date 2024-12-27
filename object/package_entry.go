package object

type PackageEntry struct {
	Hash     []byte
	Content  []byte
	Ext      string
	Filename string
}

func NewPackageEntry(hash []byte, content []byte, ext string, filename string) *PackageEntry {
	return &PackageEntry{
		Hash:     hash,
		Content:  content,
		Ext:      ext,
		Filename: filename,
	}
}
