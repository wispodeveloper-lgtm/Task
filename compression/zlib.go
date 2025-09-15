package compression

import (
	"compress/zlib"
	"fmt"
	"io"
	"os"
)

func ZlStore(inputFile string, outputFile string) (string, error) {

	// Open the input text file
	inFile, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	defer inFile.Close()

	// Create the output compressed file
	outFile, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	// Create a zlib writer
	zw := zlib.NewWriter(outFile)

	// Copy the input file into the zlib writer
	_, err = io.Copy(zw, inFile)
	if err != nil {
		panic(err)
	}

	// Close the zlib writer to flush remaining data
	if err := zw.Close(); err != nil {
		panic(err)
	}

	fmt.Println("Compression complete -> output.zlib")
	return outputFile, nil
}
