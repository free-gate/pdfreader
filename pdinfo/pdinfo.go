package main

import (
	"flag"
	"fmt"

	"github.com/raff/pdfreader/pdfread"
	"github.com/raff/pdfreader/util"
)

var (
	debugobj = false
)

// The program takes a PDF file as argument and recursively dump the content.

func printobj(pd *pdfread.PdfReaderT, o []byte, indent, prefix string) {
	l := len(o)

	if l > 2 {
		if o[l-2] == ' ' && o[l-1] == 'R' { // reference
			ref := fmt.Sprintf("<<%s>>", o)
			o = pd.Obj(o)

			if prefix == "" {
				prefix = ref
			} else {
				prefix += " " + ref
			}
		}
	}

	if l == 0 {
		fmt.Printf("%s%s %s\n", indent, prefix, "<<empty>>")
		return
	}

	if debugobj {
		fmt.Printf("%q\n", o)
	}

	switch {
	case o[0] == '[': // array
		a := pdfread.Array(o)

		fmt.Printf("%s%s %s\n", indent, prefix, "[")
		indent += "  "

		for i, v := range a {
			printobj(pd, v, indent, fmt.Sprintf("%d:", i))
		}

		indent = indent[2:]
		fmt.Printf("%s]\n", indent)

	case o[0] == '<' && o[1] == '<': // dictionary
		d := pdfread.Dictionary(o)

		fmt.Printf("%s%s %s\n", indent, prefix, "{")
		indent += "  "

		for k, v := range d {
			if k == "/Parent" { // backreference - don't follow
				fmt.Printf("%s%s <<%s>>\n", indent, k, string(v))
			} else {
				printobj(pd, v, indent, k)
			}
		}

		indent = indent[2:]
		fmt.Printf("%s}\n", indent)

	default:
		fmt.Printf("%s%s %s\n", indent, prefix, string(o))
	}
}

func main() {
	flag.BoolVar(&util.Debug, "debug", false, "enable debug logging")
	flag.BoolVar(&debugobj, "dump", false, "dump object content")

	flag.Parse()

	for _, f := range flag.Args() {
		fmt.Println("----", f, "--------------------")

		pd := pdfread.Load(f)
		if pd == nil {
			fmt.Println("can't open input file")
			fmt.Println()
			continue
		}

		fmt.Println("Trailer {")
		for k, v := range pd.Trailer {
			fmt.Printf("  %s %q\n", k, v)
		}
		fmt.Println("}")
		fmt.Println()

		root := pd.Trailer["/Root"]

		printobj(pd, root, "", "/Root")
		fmt.Println()
	}
}
