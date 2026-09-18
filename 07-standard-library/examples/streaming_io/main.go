package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

func main() {
	fmt.Println("=== 1. io.Copy and in-memory buffers ===")
	source := strings.NewReader("Hello, Gopher streaming world!\n")
	var destination bytes.Buffer

	// Copy from Reader to Writer with fixed-size internal chunks
	copied, err := io.Copy(&destination, source)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Copied %d bytes: %s", copied, destination.String())

	fmt.Println("=== 2. io.TeeReader (T-Split Stream) ===")
	input := strings.NewReader("Data to duplicate")
	var copyBuffer bytes.Buffer

	// Reading from tee simultaneously writes a duplicate to copyBuffer
	tee := io.TeeReader(input, &copyBuffer)
	mainOutput, _ := io.ReadAll(tee)

	fmt.Printf("Main Output: %s\n", string(mainOutput))
	fmt.Printf("Tee Copy:    %s\n", copyBuffer.String())

	fmt.Println("\n=== 3. io.LimitReader (Protection from DOS) ===")
	infiniteStream := strings.NewReader("MassivePayloadExceedingLimit")
	limited := io.LimitReader(infiniteStream, 7)
	limitedBytes, _ := io.ReadAll(limited)
	fmt.Printf("Read up to limit: %q\n", string(limitedBytes))
}
