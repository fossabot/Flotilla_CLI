/*
* @Author: Ximidar
* @Date:   2018-06-16 16:53:05
* @Last Modified by:   Ximidar
* @Last Modified time: 2018-06-16 17:05:26
*/

package cmd

import(
	"fmt"
	"os"
	"github.com/spf13/cobra"

)

var rootCmd = &cobra.Command{
  Use:   "mango_cli",
  Short: "Mango_cli is the cli tool for MangoOS",
  Long: `Use this tool to control MangoOS from the command line
  		 This tool will help you print files, check the status of
  		 a print, or help you control and monitor the printer command line`,
  Run: func(cmd *cobra.Command, args []string) {
    // Do Stuff Here
  },
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}
