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

/*
But the second reason is…. if you are using Cygwin/mintty/git-bash on Windows, those Windows shells are unable to reach down to the OS API, and will throw the exact same error of the “handle is invalid”.

This issue is not directly fixable and really not an issue with Go. If you switch to Powershell or even CMD then executing ReadPassword will work as expected. You may then switch back to your shell for all other commands that don’t invoke ReadPassword. If you must stay in Cygwin/minty/git-bash then take a look at https://github.com/rprichard/winpty project, it might solve your issue
*/

import (
    "os"
    "syscall"
    "golang.org/x/term"
	"log"
)

func askPassword() string {
	log.Println("Enter your Kirinuki passphrase")
	data, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Panic(err)
		os.Exit(1)
	}
	return string(data)
}