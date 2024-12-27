package pcktool

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/gopi-frame/contract/console"
	"io"
	"os"
	"pcktool/crypt"
	"pcktool/object"
	"pcktool/yappy"
	"strings"
)

func Unpack(path string, output console.Output) (*object.Package, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewReader(content)
	head, err := ReadN(buf, 8)
	fileCount, err := ReadInt(buf)
	if err != nil {
		return nil, err
	}
	pck := object.NewPackage(path, head, fileCount)
	output.Infof("Found %d files | %s", fileCount, strings.ToUpper(hex.EncodeToString(head)))
	for i := 0; i < fileCount; i++ {
		hash, err := ReadN(buf, 8)
		if err != nil {
			return nil, err
		}
		flag, err := ReadByte(buf)
		if err != nil {
			return nil, err
		}
		offset, err := ReadInt(buf)
		if err != nil {
			return nil, err
		}
		size, err := ReadInt(buf)
		if err != nil {
			return nil, err
		}
		originalSize, err := ReadInt(buf)
		if err != nil {
			return nil, err
		}
		less, err := ReadInt(buf)
		if err != nil {
			return nil, err
		}
		start, err := buf.Seek(0, io.SeekCurrent)
		if err != nil {
			return nil, err
		}
		if _, err := buf.Seek(int64(offset), io.SeekStart); err != nil {
			return nil, err
		}
		data, err := ReadN(buf, size)
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
		output.Infof("File %02d/%d: [%016X | %6d bytes or %6d] %s %02d %d",
			i+1, fileCount, offset, originalSize, size, strings.ToUpper(hex.EncodeToString(hash)), flag, less)
		pck.Entries = append(pck.Entries, entry)
		if _, err := buf.Seek(start, io.SeekStart); err != nil {
			return nil, err
		}
	}
	return pck, nil
}
