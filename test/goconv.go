package main

// iconv analog with Go
import (
	"charenc"
	"os"
	"io"
	"flag"
	"strings"
	"sort"
)

func get_locale() string {
	return "UTF-8" // TODO :)
}

func main() {
	// Parse command line:
	locale := get_locale()
	var from_enc string
	flag.StringVar(&from_enc, "from-code", locale, "convert characters from encoding.")
	flag.StringVar(&from_enc, "f", locale, "convert characters from encoding (short version).")
	var to_enc string
	flag.StringVar(&to_enc, "to-code", locale, "convert characters to encoding. If not specified the encoding corresponding to current locale is used")
	flag.StringVar(&to_enc, "t", locale, "convert characters to encoding (short version). ")
	var list bool
	flag.BoolVar(&list, "list", false, "list known code character sets.")
	flag.BoolVar(&list, "l", false, "list known code character sets (short version).")
	var output string
	flag.StringVar(&output, "output", "", "specify output file (default is stdout).")
	flag.StringVar(&output, "o", "", "specify output file (default is stdout) (short version).")
	flag.Parse()

	inputs := flag.Args()

	if list {
		array := charenc.ListEncodings()
		sort.Sort(sort.StringSlice(array))
		col := 0
		for i := range(array) {
			if i > 0 {
				print(", ")
				col += 2
			}
			if col + len(array[i]) > 72 {
				println("")
				col = 0
			}
			print(strings.ToUpper(array[i]))
			col += len(array[i])
		}

		println("")
		os.Exit(0)
	}

	var stdin io.Reader
	if len(inputs) == 0 {
		stdin = os.Stdin
	} else if len(inputs) == 1 {
		var e error
		stdin, e = os.Open(inputs[0])
		if e != nil {
			os.Stderr.Write([]byte("Error: can not open file " + inputs[0] + ":" + e.Error()))
			os.Exit(1)
		}
	} else {
		files := make([]io.Reader, len(inputs))
		var e error
		for i := range(inputs) {
			files[i], e = os.Open(inputs[i])
			if e != nil {
				os.Stderr.Write([]byte("Error: can not open file " + inputs[i] + ":" + e.Error()))
				os.Exit(1)
			}
		}

		stdin = io.MultiReader(files...)
	}

	var stdout io.WriteCloser
	if output == "" {
		stdout = os.Stdout
	} else {
		var e error
		stdout, e = os.Create(output)
		if e != nil {
			os.Stderr.Write(([]byte)("Error: can not open file " + output + ": " + e.Error()))
			os.Exit(1)
		}
	}

	reader := charenc.GetReader(stdin, from_enc, to_enc, charenc.ReplaceErrors)
	if reader == nil {
		println("Error: can not create converter")
		os.Exit(1)
	}

	buf := make([]byte, 256)
	for {
		cnt, e := reader.Read(buf)
		if cnt > 0 {
			cnt2, e2 := stdout.Write(buf)
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

	stdout.Close()
}

