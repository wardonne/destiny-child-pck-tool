package yappy

var infos = make([]int, 256)
var maps = make([][]int, 32)

func init() {
	for i := range maps {
		maps[i] = make([]int, 16)
	}
	step := 1 << 16
	for i := 0; i < 16; i++ {
		value := 65535
		step = (step * 67537) >> 16
		for value < (29 << 16) {
			maps[value>>16][i] = 1
			value = (value * step) >> 16
		}
	}
	c := 0
	for i := 0; i < 29; i++ {
		for j := 0; j < 16; j++ {
			if maps[i][j] == 1 {
				infos[32+c] = i + 4 + (j << 8)
				maps[i][j] = 32 + c
				c++
			} else {
				if i == 0 {
					panic("i == 0")
				}
				maps[i][j] = maps[i-1][j]
			}
		}
	}
	if c != 256-32 {
		panic("init error")
	}
}

func Decompress(data []byte, size int) ([]byte, error) {
	to := make([]byte, 0, size)
	dataP := 0
	toP := 0
	for len(to) < size {
		if !(dataP+1 < len(data)) {
			return data, nil
		}
		index := data[dataP] & 0xFF
		if index < 32 {
			to = append(to, data[dataP+1:dataP+1+int(index)+1]...)
			toP += int(index) + 1
			dataP += int(index) + 2
		} else {
			info := infos[index]
			length := info & 0xFF
			offset := (info & 0xFF00) + int(data[dataP+1]&0xFF)
			to = append(to, to[toP-offset:toP-offset+length]...)
			toP += length
			dataP += 2
		}
	}
	return to, nil
}
