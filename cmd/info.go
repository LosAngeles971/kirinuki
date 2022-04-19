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
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "information of file",
	Long: `information of file.
Usage:
	kirinuki info --name <name>`,
	Run: func(cmd *cobra.Command, args []string) {
		g := getGateway(email, askPassword())
		err := g.Login()
		if err != nil {
			log.Fatalf("login failed [%v]", err)
		}
		rr, err := g.Find(name)
		if err != nil {
			log.Fatal(err)
		}
		for _, k := range rr {
			log.Infof("kirinuki %s - date %v ", k.Name, k.Date)
		}
	},
}

func init() {
	infoCmd.PersistentFlags().StringVar(&name, "name", "", "pattern for name finding")
	infoCmd.MarkPersistentFlagRequired("name")
	rootCmd.AddCommand(infoCmd)
}


