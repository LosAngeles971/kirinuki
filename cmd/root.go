/*
 * Created on Sun Apr 10 2022
 * Author @LosAngeles971
 *
 * This software is licensed under GNU General Public License v2.0
 * Copyright (c) 2022 @LosAngeles971
 *
 * The GNU GPL is the most widely used free software license and has a strong copyleft requirement.
 * When distributing derived works, the source code of the work must be made available under the same license.
 * There are multiple variants of the GNU GPL, each with different requirements.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED
 * TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 * THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
 * TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */
package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var name string
var filename string
var email string
var overwrite bool
var debug bool
var storageMap string

func setLog(cmd *cobra.Command, args []string) {
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("debug level activated")
	}
}

var rootCmd = &cobra.Command{
	Use:   "kirinuki",
	Short: "Kirinuki",
	Long:  `Kirinuki.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
	},
	PersistentPreRun: setLog,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&email, "email", "", "email")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug")
	rootCmd.PersistentFlags().StringVar(&storageMap, "storage", "", "storage file config")
	rootCmd.MarkPersistentFlagRequired("email")
}
