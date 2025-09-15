package compression

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

func GzStore(inputFile, outputFile string) (string, error) {

	// Open the source file
	inFile, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	defer inFile.Close()

	// Create the destination file
	outFile, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	// Create a gzip writer on top of the output file
	gzipWriter := gzip.NewWriter(outFile)
	defer gzipWriter.Close()

	// Copy data from the input file to the gzip writer
	_, err = io.Copy(gzipWriter, inFile)
	if err != nil {
		panic(err)
	}

	fmt.Println("File compressed successfully:", outputFile)
	return outputFile, nil
}
