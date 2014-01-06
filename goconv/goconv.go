package main

// iconv analog with Go
import (
	"charenc"
	"os"
	"io"
	"flag"
	"strings"
	"sort"
	"fmt"
)

func get_locale() string {
	return "UTF-8" // TODO :)
}

type cmdline struct {
	from_enc, to_enc string
	list bool
	output string
	flags int
	inputs []string
}

func parse_cmdline() cmdline {
	var r cmdline
	var replace, ignore bool

	locale := get_locale()
	flag.StringVar(&r.from_enc, "from-code", locale, "convert characters from encoding.")
	flag.StringVar(&r.from_enc, "f", locale, "convert characters from encoding (short version).")
	flag.StringVar(&r.to_enc, "to-code", locale, "convert characters to encoding. If not specified the encoding corresponding to current locale is used")
	flag.StringVar(&r.to_enc, "t", locale, "convert characters to encoding (short version). ")
	flag.BoolVar(&r.list, "list", false, "list known code character sets.")
	flag.BoolVar(&r.list, "l", false, "list known code character sets (short version).")
	flag.StringVar(&r.output, "output", "", "specify output file (default is stdout).")
	flag.StringVar(&r.output, "o", "", "specify output file (default is stdout) (short version).")
	flag.BoolVar(&replace, "r", false, "replace invalid characters in input and output streams")
	flag.BoolVar(&ignore, "c", false, "ignore invalid characters in input and output streams")
	var help bool
	flag.BoolVar(&help, "help", false, "print this help and exit")
	flag.BoolVar(&help, "h", false, "print this help and exit (short version)")

	flag.Parse()

	if help {
		fmt.Fprintf(os.Stderr, "Usage: %s [-f encoding] [-t encoding] [inputfile]...\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}

	r.inputs = flag.Args()

	return r
}

func print_list() {
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

func main() {
	// Parse command line:
	params := parse_cmdline()

	if params.list {
		print_list()
	}

	var stdin io.Reader
	if len(params.inputs) == 0 {
		stdin = os.Stdin
	} else if len(params.inputs) == 1 {
		var e error
		stdin, e = os.Open(params.inputs[0])
		if e != nil {
			fmt.Fprintf(os.Stderr, "Error: can not open file '%s': %s\n", params.inputs[0], e.Error())
			os.Exit(1)
		}
	} else {
		files := make([]io.Reader, len(params.inputs))
		var e error
		for i := range(params.inputs) {
			files[i], e = os.Open(params.inputs[i])
			if e != nil {
				fmt.Fprintf(os.Stderr, "Error: can not open file '%s': %s\n", params.inputs[i], e.Error())
				os.Exit(1)
			}
		}

		stdin = io.MultiReader(files...)
	}

	var stdout io.WriteCloser
	if params.output == "" {
		stdout = os.Stdout
	} else {
		var e error
		stdout, e = os.Create(params.output)
		if e != nil {
			fmt.Fprintf(os.Stderr, "Error: can not open file '%s': %s\n", params.output, e.Error())
			os.Exit(1)
		}
	}

	reader := charenc.GetReader(stdin, params.from_enc, params.to_enc, charenc.ReplaceErrors)
	if reader == nil {
		fmt.Fprintf(os.Stderr, "Error: can not create converter\n")
		os.Exit(1)
	}

	buf := make([]byte, 256)
	for {
		cnt, e := reader.Read(buf)
		if cnt > 0 {
			cnt2, e2 := stdout.Write(buf[:cnt])
			if cnt2 < cnt {
				if e2 != nil { // MUST work here
					fmt.Fprintf(os.Stderr, "Write failed: %s\n", e.Error())
				} else {
					fmt.Fprintf(os.Stderr, "Write failed: unknown error\n")
				}
				os.Exit(1)
			}
		}
		if e == io.EOF {
			break
		} else if e != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", e.Error())
			os.Exit(1)
		}
	}

	stdout.Close()
}

