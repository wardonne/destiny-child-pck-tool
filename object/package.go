package object

type Package struct {
	Path      string
	Head      []byte
	FileCount int
	Entries   []*PackageEntry
}

func NewPackage(path string, head []byte, fileCount int) *Package {
	return &Package{
		Path:      path,
		Head:      head,
		FileCount: fileCount,
		Entries:   make([]*PackageEntry, 0),
	}
}
