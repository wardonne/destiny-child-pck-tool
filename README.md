# Destiny Child Pck Tool

This is a package for unpacking Pck files from Destiny Child (kr).

## Import
```shell
go get -u github.com/wardonne/pck
```

## Usage

### Unpack a Pck file
```golang
package main

import pcktool "github.com/wardonne/destiny-child-pck-tool"
import "fmt"

func main() {
	pck, err := pcktool.Unpack("input.pck")
	if err != nil {
		panic(err)
	}
	for _, file := range pck.Entries {
		// do something with file
    }
}
```

### Unpack a Live2D Pck file
```go
package main

import pcktool "github.com/wardonne/destiny-child-pck-tool"
import "fmt"

func main() {
	pck, err := pcktool.Unpack("input.pck")
	if err != nil {
		panic(err)
	}
	model, err := pcktool.GenerateLive2D(pck)
	if err != nil {
		panic(err)
	}
	fmt.Println(model)
	for _, file := range pck.Entries {
		// do something with file
	}
}
```

## Command Line Tool

This package also provides a simple command line tool for unpacking Pck files.

```shell
# install the tool
go install "github.com/wardonne/destiny-child-pck-tool/cmd/pcktool"

# unpack pck file(s)
pcktool -s <source directory or file> -t <output directory>

# unpack pck file(s) as live2d
pcktool -s <source directory or file> -t <output directory> -l

# for more information
pcktool -h
```
