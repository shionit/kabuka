/*
Copyright © 2021 shionit

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/shionit/kabuka/internal/app/kabuka"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kabuka",
	Short: "Show stock information",
	Args:  cobra.MinimumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		symbol := cmd.Flags().Arg(0) // Ticker like "3994.T"

		kabuka := &kabuka.Kabuka{
			Option: kabuka.Option{
				Symbol: symbol,
			},
		}
		result, err := kabuka.Fetch()
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("%s kabuka is %s\n", result.Symbol, result.CurrentPrice)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kabuka.yaml)")
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
