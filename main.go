/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"flag"

	"github.com/cldmnky/exsconvert/cmd"
	"k8s.io/klog"
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()
	cmd.Execute()
}
