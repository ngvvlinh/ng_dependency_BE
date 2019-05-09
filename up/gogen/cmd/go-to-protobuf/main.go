package main

import (
	goflag "flag"

	flag "github.com/spf13/pflag"
	"k8s.io/code-generator/cmd/go-to-protobuf/protobuf"
	"k8s.io/klog"
)

var g = protobuf.New()

func init() {
	klog.InitFlags(goflag.CommandLine)

	g.BindFlags(flag.CommandLine)
	goflag.Set("logtostderr", "true")
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
}

func main() {
	flag.Parse()

	g.APIMachineryPackages = ""
	protobuf.Run(g)
}
