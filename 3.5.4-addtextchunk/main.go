package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"math"
	"os"
)

func textChunk(t string) (io.Reader, error) {
	b := []byte(t)
	if len(b) > math.MaxInt32 {
		return nil, fmt.Errorf("text size error: len(text)=%v > MaxInt32=%v", len(b), math.MaxInt32)
	}
	buff := &bytes.Buffer{}

	if err := binary.Write(buff, binary.BigEndian, int32(len(b))); err != nil {
		return nil, fmt.Errorf("Write error: %w", err)
	}
	if _, err := buff.WriteString("tEXt"); err != nil {
		return nil, fmt.Errorf("WriteString error: %w", err)
	}
	if _, err := buff.Write(b); err != nil {
		return nil, fmt.Errorf("Write error: %w", err)
	}

	// CRC を計算して追加
	crc := crc32.NewIEEE()
	if _, err := io.WriteString(crc, "tEXt"); err != nil {
		return nil, fmt.Errorf("WriteString error: %w", err)
	}
	if _, err := crc.Write(b); err != nil {
		return nil, fmt.Errorf("Write error: %w", err)
	}
	if err := binary.Write(buff, binary.BigEndian, crc.Sum32()); err != nil {
		return nil, fmt.Errorf("Write error: %w", err)
	}
	return buff, nil
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

func copyPNGWithTextChunk(src, dest, text string) error {
	f, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("Open error: %w", err)
	}
	defer f.Close()
	cs, err := readChunks(f)
	if err != nil {
		return fmt.Errorf("readChunks error: %w", err)
	}
	if len(cs) < 2 {
		return fmt.Errorf("chunks size error: len(chunks)=%v", len(cs))
	}

	buff := &bytes.Buffer{}
	// シグニチャの書き込み
	if _, err := io.WriteString(buff, "\x89PNG\r\n\x1a\n"); err != nil {
		return fmt.Errorf("WriteString error: %w", err)
	}
	// 先頭に必要な IHDR チャンクを書き込み
	if _, err := io.Copy(buff, cs[0]); err != nil {
		return fmt.Errorf("Copy error: %w", err)
	}
	// テキストチャンクを追加
	tc, err := textChunk("I LOVE DOG")
	if err != nil {
		return fmt.Errorf("textChunk error: %w", err)
	}
	if _, err := io.Copy(buff, tc); err != nil {
		return fmt.Errorf("Copy error: %w", err)
	}
	// 残りのチャンクを追加
	for _, c := range cs[1:] {
		if _, err := io.Copy(buff, c); err != nil {
			return fmt.Errorf("Copy error: %w", err)
		}
	}

	// チャンク追加が成功したらファイルに保存
	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("Create error: %w", err)
	}
	defer out.Close()
	if _, err := io.Copy(out, buff); err != nil {
		return fmt.Errorf("Copy error: %w", err)
	}

	return nil
}

func main() {
	err := copyPNGWithTextChunk("dog.png", "dogWithTextChunk.png", "I LOVE DOG")
	if err != nil {
		panic(err)
	}
}
