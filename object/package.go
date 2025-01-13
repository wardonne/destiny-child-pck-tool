package object

type Package struct {
	Path      string          `json:"path,omitempty"`
	Head      []byte          `json:"head,omitempty"`
	FileCount int             `json:"file_count,omitempty"`
	Entries   []*PackageEntry `json:"entries,omitempty"`
}

func NewPackage(path string, head []byte, fileCount int) *Package {
	return &Package{
		Path:      path,
		Head:      head,
		FileCount: fileCount,
		Entries:   make([]*PackageEntry, 0),
	}
}
