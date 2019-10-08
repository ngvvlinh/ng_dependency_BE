package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"etop.vn/backend/tools/pkg/generator"
	"etop.vn/backend/tools/pkg/generators/convert"
	"etop.vn/backend/tools/pkg/generators/cq"
	"etop.vn/backend/tools/pkg/generators/event"
	"etop.vn/backend/tools/pkg/generators/sample"
)

var flClean = flag.Bool("clean", false, "clean generated files without generating new files")
var flPlugins = flag.String("plugins", "", "comma separated list of plugins for generating (default to all plugins)")
var flIgnoredPlugins = flag.String("ignored-plugins", "", "comma separated list of plugins to ignore")

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
		event.New(),
		sample.New(),
	}

	flag.Parse()
	patterns := flag.Args()
	if len(patterns) == 0 {
		usage()
		os.Exit(2)
	}

	enabledPlugins := allPluginNames(plugins)
	if *flPlugins != "" {
		enabledPlugins = strings.Split(*flPlugins, ",")
	}
	if *flIgnoredPlugins != "" {
		ignoredPlugins := strings.Split(*flIgnoredPlugins, ",")
		enabledPlugins = calcEnabledPlugins(enabledPlugins, ignoredPlugins)
	}

	cfg := generator.Config{
		CleanOnly:      *flClean,
		Namespace:      "etop.vn",
		EnabledPlugins: enabledPlugins,
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

func allPluginNames(plugins []generator.Plugin) []string {
	names := make([]string, len(plugins))
	for i, p := range plugins {
		names[i] = p.Name()
	}
	return names
}

func calcEnabledPlugins(plugins []string, ignoredPlugins []string) []string {
	result := make([]string, 0, len(plugins))
	for _, p := range plugins {
		include := true
		for _, ip := range ignoredPlugins {
			if p == ip {
				include = false
				break
			}
		}
		if include {
			result = append(result, p)
		}
	}
	return result
}
