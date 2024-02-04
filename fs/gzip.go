package fs

import (
	"compress/gzip"
	"io"
	"os"
)

func CompressFile(inputPath string, outputPath string) error {
	inFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inFile.Close()
	// 创建gzip压缩输出流
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// 创建gzip.Writer包裹输出文件
	gzw := gzip.NewWriter(outFile)
	defer gzw.Close()

	// 将输入文件内容复制到gzip写入器
	_, err = io.Copy(gzw, inFile)
	if err != nil {
		return err
	}

	return nil
}

func DecompressFile(inputPath string, outputPath string) error {
	inFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	gzr, err := gzip.NewReader(inFile)
	if err != nil {
		return err
	}
	defer gzr.Close()

	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, gzr)
	if err != nil {
		return err
	}

	return nil
}
