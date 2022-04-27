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
	"time"

	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var findCmd = &cobra.Command{
	Use:   "find",
	Short: "find files into the storage using regex matching",
	Long: `find files into the storage regex matching.
Usage:
	kirinuki find --name <regex pattern>`,
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
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Date"})
		for _, k := range rr {
			table.Append([]string{
				k.Name,
				time.Unix(k.Date, 0).String(),
			})		
		}
		table.Render()
	},
}

func init() {
	findCmd.PersistentFlags().StringVar(&name, "name", "", "pattern for name finding")
	findCmd.MarkPersistentFlagRequired("name")
	rootCmd.AddCommand(findCmd)
}


