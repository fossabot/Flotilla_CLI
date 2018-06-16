/*
* @Author: Ximidar
* @Date:   2018-06-16 16:39:58
* @Last Modified by:   Ximidar
* @Last Modified time: 2018-06-16 16:54:19
*/

package cmd

import(
	"fmt"
	"github.com/spf13/cobra"

)

func init() {
  rootCmd.AddCommand(printerface)
}

var printerface = &cobra.Command{
  Use:   "ui",
  Short: "Show the cli UI for MangoOS",
  Long:  `This will open the cli UI for MangoOS. This has tools for monitoring the command line and starting prints (or it will in the future)`,
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("Hello, I work!")
  },
}

