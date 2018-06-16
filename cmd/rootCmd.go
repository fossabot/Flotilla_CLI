/*
* @Author: Ximidar
* @Date:   2018-06-16 16:53:05
* @Last Modified by:   Ximidar
* @Last Modified time: 2018-06-16 16:55:42
*/

package cmd

import(
	"fmt"
	"os"
	"github.com/spf13/cobra"

)

var rootCmd = &cobra.Command{
  Use:   "mango_cli",
  Short: "MangoCLI is the cli tool for MangoOS",
  Long: `Fill this out later`,
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
