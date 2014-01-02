package main

// iconv analog with Go
import (
	"charenc"
	"os"
	"io"
)

func main() {
	reader := charenc.GetReader(os.Stdin, "koi8-r", "utf-8", charenc.ReplaceErrors)
	if reader == nil {
		println("Error: can not create converter")
		os.Exit(1)
	}

	writer := os.Stdout
	buf := make([]byte, 256)
	for {
		cnt, e := reader.Read(buf)
		if cnt > 0 {
			cnt2, e2 := writer.Write(buf)
			if cnt2 < cnt {
				if e2 != nil { // MUST be
					println("Write failed: " + e.Error())
				} else {
					println("Write failed: unknown error")
				}
				os.Exit(1)
			}
		}
		if e == io.EOF {
			break
		} else if e != nil {
			println("Error: " + e.Error())
			os.Exit(1)
		}
	}
}

