/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"mytools/tools/recactibw"
	"mytools/utils"
)

var (
	source            string
	percent           float64
	reCactiBwInterval int64
	start             string
	end               string
	workerCount       int
	details           bool
	logFile           string
	print             bool
)

// recactibwCmd represents the recactibw command
var recactibwCmd = &cobra.Command{
	Use:   "recactibw",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("recactibw called")
		e := recactibw.Engine{
			Scheduler:   &recactibw.QueueScheduler{},
			WorkerCount: workerCount,
		}
		e.Run(utils.StringToTimestamp(start), utils.StringToTimestamp(end),
			reCactiBwInterval*86400, percent, source, details, logFile, print)
	},
}

func init() {
	rootCmd.AddCommand(recactibwCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// recactibwCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// recactibwCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	recactibwCmd.Flags().StringVarP(&source, "source", "s", "", "需要更新的目标文件或文件夹")
	recactibwCmd.Flags().Float64VarP(&percent, "percent", "p", 1, "将一定范围内的数据按照百分百升降.0.3(表示减小到30%)")
	recactibwCmd.Flags().Int64VarP(&reCactiBwInterval, "interval", "i", 0, "参照若干时间点前的数据修补当前数据单位为天")
	recactibwCmd.Flags().StringVarP(&start, "start", "", "", "修改起始时间点2020-01-21 00:00:00")
	recactibwCmd.Flags().StringVarP(&end, "end", "", "", "修改起始时间点2020-01-21 00:00:00")
	recactibwCmd.Flags().IntVarP(&workerCount, "workercount", "w", 5, "同时运行修改文件的数量")
	recactibwCmd.Flags().BoolVarP(&details, "details", "d", false, "是否打印详情")
	recactibwCmd.Flags().StringVarP(&logFile, "log", "l", "/var/log/mytools/recactibw.log", "操作文件日志所存地方")
	recactibwCmd.Flags().BoolVarP(&print, "print", "", false, "显示数值模式")

}
