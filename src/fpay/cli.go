/* The MIT License (MIT)
Copyright © 2018 by Atlas Lee(atlas@fpay.io)

Permission is hereby granted, free of charge, to any person obtaining a
copy of this software and associated documentation files (the “Software”),
to deal in the Software without restriction, including without limitation
the rights to use, copy, modify, merge, publish, distribute, sublicense,
and/or sell copies of the Software, and to permit persons to whom the
Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
DEALINGS IN THE SOFTWARE.
*/

package fpay

import (
	"flag"
)

type Settings struct {
	Anumber                            int
	Laddr, Raddr, Gblock, Apath, Paddr string
	Naccounts                          bool
}

// 命令行参数解释
func ParseCLI() (settings *Settings, err error) {
	settings = new(Settings)
	flag.StringVar(&settings.Laddr, "L", ":8080", "Listening address. Default is :8080")
	flag.StringVar(&settings.Raddr, "R", ":6379", "Redis IP Address. Default is :6379.")
	flag.StringVar(&settings.Gblock, "G", "", "God block config.")
	flag.StringVar(&settings.Apath, "A", "", "Account file.")
	flag.StringVar(&settings.Paddr, "P", "", "Default parent address.")
	flag.BoolVar(&settings.Naccounts, "N", false, "Generate new accounts.")
	flag.IntVar(&settings.Anumber, "AN", 1, "Number of new accounts generated.")
	flag.Parse()
	return
}
