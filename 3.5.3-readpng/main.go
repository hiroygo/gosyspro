package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func dumpChunk(c io.Reader) error {
	var length int32
	// 普通 io.Reader は Read で Seek するはずなので
	// Read する度に読み取るバイト位置が移動していく
	err := binary.Read(c, binary.BigEndian, &length)
	if err != nil {
		return fmt.Errorf("Read error: %w", err)
	}

	buff := make([]byte, 4)
	// ReadAt は 4 バイト読み取れないとエラーを返すけど
	// Read はエラーを返さないので、自分で読み取ったバイト数を調べる必要がある
	n, err := c.Read(buff)
	if err != nil {
		return fmt.Errorf("Read error: %w", err)
	}
	if n != len(buff) {
		return fmt.Errorf("Read error: got %v bytes", n)
	}

	fmt.Printf("chunk %q (%v bytes)\n", string(buff), length)

	if bytes.Equal(buff, []byte("tEXt")) {
		rawText := make([]byte, length)
		if _, err := c.Read(rawText); err != nil {
			return fmt.Errorf("Read error: %w", err)
		}
		fmt.Println(string(rawText))
	}

	return nil
}

func readChunks(f *os.File) ([]io.Reader, error) {
	if _, err := f.Seek(8, io.SeekStart); err != nil {
		return nil, fmt.Errorf("Seek error: %w", err)
	}

	rs := make([]io.Reader, 0)
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

		// NewSectionReader では f が既に Seek されていても offset から Read される
		rs = append(rs, io.NewSectionReader(f, offset, int64(length+12)))

		// 現在位置は長さを読み終わった箇所なので、次のチャンクの先頭に移動
		offset, err = f.Seek(int64(length+8), io.SeekCurrent)
		if err != nil {
			return nil, fmt.Errorf("Seek error: %w", err)
		}
	}
	return rs, nil
}

func main() {
	if len(os.Args) != 2 {
		panic("len(Args) != 2")
	}

	f, err := os.Open(os.Args[1])
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
