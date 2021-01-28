/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/spf13/cobra"
	"mytools/tools/mping"
)

var (
	count       int
	spacing     int
	workercount int
	repeat      int
	interval    int
	ipfile      string
	ipnet       string
)

// mpingCmd represents the mping command
var mpingCmd = &cobra.Command{
	Use:   "mping",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		pingers := mping.NewPingers(count, spacing, workercount, repeat, interval, ipfile, ipnet)
		pingers.Run()
	},
}

func init() {
	rootCmd.AddCommand(mpingCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mpingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mpingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	mpingCmd.Flags().IntVarP(&count, "count", "c", 5, "每次每个对象测试包数")
	mpingCmd.Flags().IntVarP(&spacing, "spacing", "s", 2, "每次测试时间间隔")
	mpingCmd.Flags().IntVarP(&workercount, "workercount", "w", 4, "同时测试的协程数量")
	mpingCmd.Flags().IntVarP(&repeat, "repeat", "r", 1, "所有IP测试的遍数，默认1次")
	mpingCmd.Flags().IntVarP(&interval, "interval", "i", 500, "每个测试对象发包的时间间隔，单位毫秒")
	mpingCmd.Flags().StringVarP(&ipfile, "ipfile", "f", "", "从文件中获取测试IP地址")
	mpingCmd.Flags().StringVarP(&ipnet, "ipnet", "n", "", "根据网段获取测试IP地址，例如192.168.0.0/24")
}
