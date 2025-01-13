package pcktool

import (
	"bytes"
	"fmt"
	"github.com/wardonne/destiny-child-pck-tool/crypt"
	"github.com/wardonne/destiny-child-pck-tool/object"
	"github.com/wardonne/destiny-child-pck-tool/yappy"
	"io"
	"os"
)

func Unpack(path string) (*object.Package, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewReader(content)
	pkg, err := UnpackStream(buf)
	if err != nil {
		return nil, err
	}
	pkg.Path = path
	return pkg, nil
}

func UnpackStream(stream io.ReadSeeker) (*object.Package, error) {
	head, err := getHead(stream)
	if err != nil {
		return nil, err
	}
	fileCount, err := getFileCount(stream)
	if err != nil {
		return nil, err
	}
	pck := object.NewPackage("", head, fileCount)
	entries, err := getFiles(stream, fileCount)
	if err != nil {
		return nil, err
	}
	pck.Entries = entries
	return pck, nil
}

func getHead(stream io.Reader) ([]byte, error) {
	head, err := ReadN(stream, 8)
	if err != nil {
		return nil, err
	}
	return head, nil
}

func getFileCount(stream io.Reader) (int, error) {
	fileCount, err := ReadInt(stream)
	if err != nil {
		return 0, err
	}
	return fileCount, nil
}

func getFiles(stream io.ReadSeeker, fileCount int) ([]*object.PackageEntry, error) {
	var entries []*object.PackageEntry
	for i := 0; i < fileCount; i++ {
		hash, err := ReadN(stream, 8)
		if err != nil {
			return nil, err
		}
		flag, err := ReadByte(stream)
		if err != nil {
			return nil, err
		}
		offset, err := ReadInt(stream)
		if err != nil {
			return nil, err
		}
		size, err := ReadInt(stream)
		if err != nil {
			return nil, err
		}
		originalSize, err := ReadInt(stream)
		if err != nil {
			return nil, err
		}
		_, err = ReadInt(stream)
		if err != nil {
			return nil, err
		}
		start, err := stream.Seek(0, io.SeekCurrent)
		if err != nil {
			return nil, err
		}
		if _, err := stream.Seek(int64(offset), io.SeekStart); err != nil {
			return nil, err
		}
		data, err := ReadN(stream, size)
		if err != nil {
			return nil, err
		}
		if flag&2 == 2 {
			data, err = crypt.Decrypt(data)
			if err != nil {
				return nil, err
			}
		}
		if flag&1 == 1 {
			data, err = yappy.Decompress(data, originalSize)
			if err != nil {
				return nil, err
			}
		}
		var ext byte
		if len(data) > 0 {
			ext = data[0] & 0xFF
		}
		var extMap = map[byte]string{
			109: "dat",
			35:  "mtn",
			137: "png",
			123: "json",
		}
		extStr, ok := extMap[ext]
		if !ok {
			extStr = "unk"
		}
		entry := object.NewPackageEntry(hash, data, extStr, fmt.Sprintf("%08d.%s", i, extStr))
		if _, err := stream.Seek(start, io.SeekStart); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
		if _, err := stream.Seek(start, io.SeekStart); err != nil {
			return nil, err
		}
	}
	return entries, nil
}
