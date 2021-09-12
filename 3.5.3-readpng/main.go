package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// io.SectionReader は Read すると Seek する
func dumpChunk(c *io.SectionReader) error {
	var length int32
	err := binary.Read(c, binary.BigEndian, &length)
	if err != nil {
		return fmt.Errorf("Read error: %w", err)
	}
	buff := make([]byte, 4)
	// SectionReader 型として一度 binary.Read して
	// Seek 済なので、このコードは問題ないはず
	// c が io.Reader だと Seek が保証できない
	_, err = c.Read(buff)
	if err != nil {
		return fmt.Errorf("Read error: %w", err)
	}
	fmt.Printf("chunk %q (%v bytes)\n", string(buff), length)
	return nil
}

func readChunks(f *os.File) ([]*io.SectionReader, error) {
	if _, err := f.Seek(8, io.SeekStart); err != nil {
		return nil, fmt.Errorf("Seek error: %w", err)
	}

	rs := make([]*io.SectionReader, 0)
	var offset int64 = 8
	for {
		var length int32
		err := binary.Read(f, binary.BigEndian, &length)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Read error: %w", err)
		}

		// NewSectionReader では f が Seek されていても offset から Read される
		rs = append(rs, io.NewSectionReader(f, offset, int64(length+12)))

		// 現在位置は長さを読み終わった箇所なので
		// 次のチャンクの先頭に移動
		offset, err = f.Seek(int64(length+8), io.SeekCurrent)
	}
	return rs, nil
}

func main() {
	f, err := os.Open("dog.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	cs, err := readChunks(f)
	if err != nil {
		panic(err)
	}

	for _, c := range cs {
		err := dumpChunk(c)
		if err != nil {
			panic(err)
		}
	}
}
