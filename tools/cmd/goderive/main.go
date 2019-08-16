//  Copyright 2017 Walter Schulze
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

// Package main implements the goderive binary.
// This pulls in all the plugins, parses the flags and runs the generators using the derive library.
package main

import (
	"flag"
	"log"
	"strings"

	"github.com/awalterschulze/goderive/derive"

	"etop.vn/backend/tools/pkg/goderive/sel"
	"etop.vn/backend/tools/pkg/goderive/substruct"
	"etop.vn/backend/tools/pkg/sqlgen"
)

var autoname = flag.Bool("autoname", false, "rename functions that are conflicting with other functions")
var dedup = flag.Bool("dedup", false, "rename functions to functions that are duplicates")
var prefix = flag.String("prefix", "derive", "prefix of all functions")
var pluginprefix = flag.String("pluginprefix", "", "used to override function prefixes.  The input is a comma separated list of function are prefix pairs.  For example equal=deriveEqual,copyto=copyTo,fmap=fmap,")

func main() {
	plugins := []derive.Plugin{
		substruct.NewPlugin(),
		sqlgen.NewPlugin(),
		sel.NewPlugin(),
	}
	log.SetFlags(0)
	flag.Parse()
	overridePrefixes := make(map[string]string)
	if len(*pluginprefix) > 0 {
		pairs := strings.Split(*pluginprefix, ",")
		for _, pair := range pairs {
			ss := strings.Split(pair, "=")
			if len(ss) != 2 {
				log.Fatalf("invalid syntax for plugin prefix <%s>", pair)
			}
			overridePrefixes[ss[0]] = ss[1]
		}
	}
	for _, p := range plugins {
		pluginprefix := p.GetPrefix()
		pluginprefix = strings.Replace(pluginprefix, "derive", *prefix, 1)
		newprefix, override := overridePrefixes[p.Name()]
		if override {
			pluginprefix = newprefix
		}
		p.SetPrefix(pluginprefix)
	}
	paths := derive.ImportPaths(flag.Args())
	g, err := derive.NewPlugins(plugins, *autoname, *dedup).Load(paths)
	if err != nil {
		log.Fatal(err)
	}
	if err := g.Generate(); err != nil {
		log.Fatal(err)
	}
}
