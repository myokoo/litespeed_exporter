package main

import (
	"fmt"

	"github.com/myokoo/litespeed_exporter/pkg/rtreport"
)

func main() {
	v, err := rtreport.New(rtreport.DefaultReportPath)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	fmt.Printf("Version: %.1s\n", v.Version)
	fmt.Printf("Uptime: %.0f\n", v.Uptime)
	fmt.Println("ConnectionReport:")
	for key, value := range v.ConnectionReport {
		fmt.Printf("\t%s: %.1f\n", key, value)
	}
	fmt.Println("NetworkReport:")
	for key, value := range v.NetworkReport {
		fmt.Printf("\t%s: %.1f\n", key, value)
	}
	fmt.Println("RequestReports:")
	for key, _ := range v.RequestReports {
		fmt.Printf("\t%s:\n", key)
		for key2, value2 := range v.RequestReports[key] {
			fmt.Printf("\t\t%s: %.1f\n", key2, value2)
		}
	}
	fmt.Println("ExtAppReports:")
	for key, _ := range v.ExtAppReports {
		fmt.Printf("\t%s:\n", key)
		for key2, _ := range v.ExtAppReports[key] {
			fmt.Printf("\t\t%s:\n", key2)
			for key3, _ := range v.ExtAppReports[key][key2] {
				fmt.Printf("\t\t\t%s:\n", key3)
				for key4, value4 := range v.ExtAppReports[key][key2][key3] {
					fmt.Printf("\t\t\t\t%s:%.1f\n", key4, value4)

				}
			}
		}
	}
}
