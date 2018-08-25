/*
* @Author: Ximidar
* @Date:   2018-06-16 16:53:05
* @Last Modified by:   Ximidar
* @Last Modified time: 2018-08-25 10:09:18
*/

package cli_tools

import(
	"fmt"
	"os"
	"github.com/spf13/cobra"
  "github.com/ximidar/mango_cli/user_interface"

)

var rootCmd = &cobra.Command{
  Use:   "mango_cli",
  Short: "Mango_cli is the cli tool for MangoOS",
  Long: `Use this tool to control MangoOS from the command line
  		 This tool will help you print files, check the status of
  		 a print, or help you control and monitor the printer command line`,
  Run: func(cmd *cobra.Command, args []string) {
  	if len(args) == 0 {
  		cmd.Help()
  		os.Exit(1)
  	} else {
  		fmt.Println("Mango CLI")
  		fmt.Println("Written By: Matt Pedler")
  	}
  },
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func init() {
  rootCmd.AddCommand(printerface)
}

var printerface = &cobra.Command{
  Use:   "ui",
  Short: "Show the cli UI for MangoOS",
  Long:  `This will open the cli UI for MangoOS. This has tools for monitoring the command line and starting prints (or it will in the future)`,
  Run: func(cmd *cobra.Command, args []string) {
    cligui := user_interface.New_Cli_Gui()
    cligui.Screen_Init()
  },
}
