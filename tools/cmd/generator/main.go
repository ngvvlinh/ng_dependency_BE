package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"etop.vn/backend/tools/pkg/generator"
	"etop.vn/backend/tools/pkg/generators/convert"
	"etop.vn/backend/tools/pkg/generators/cq"
	"etop.vn/backend/tools/pkg/generators/sample"
)

var flClean = flag.Bool("clean", false, "clean generated files without generating new files")
var flPlugins = flag.String("plugins", "", "select plugins for generating (default to all plugins)")

func usage() {
	const text = `
Usage: generator [OPTION] PATTERN ...

Options:
`
	fmt.Print(text[1:])
	flag.PrintDefaults()
}

func main() {
	plugins := []generator.Plugin{
		convert.New(),
		cq.New(),
		sample.New(),
	}

	flag.Parse()
	patterns := flag.Args()
	if len(patterns) == 0 {
		usage()
		os.Exit(2)
	}

	cfg := generator.Config{
		CleanOnly: *flClean,
		Namespace: "etop.vn",
	}
	if *flPlugins != "" {
		cfg.EnabledPlugins = strings.Split(*flPlugins, ",")
	}
	if err := generator.RegisterPlugin(plugins...); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
	if err := generator.Start(cfg, patterns...); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}
