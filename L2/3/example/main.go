package main

import (
	"flag"
	"fmt"
	"os"
	"sortutil"
	"strings"
)

func init() {
	if len(os.Args) < 2 {
		os.Exit(1)
	}
	sortutil.FileNameSortOpt = os.Args[len(os.Args)-1]
	flag.BoolVar(&sortutil.DebugSortOpt, "d", false, "debug")
	flag.BoolVar(&sortutil.NumberSortOpt, "n", false, "сортировать по числовому значению")
	flag.BoolVar(&sortutil.ReverseSortOpt, "r", false, "сортировать в обратном порядке")
	flag.BoolVar(&sortutil.UniqueSortOpt, "u", false, "не выводить повторяющиеся строки")
	flag.IntVar(&sortutil.ColumnSortOpt, "k", -1, "указание колонки для сортировки (слова в строке могут выступать в качестве колонок, по умолчанию разделитель — пробел)")
	flag.StringVar(&sortutil.ColumnSeparatorSortOpt, "s", " ", "слова в строке могут выступать в качестве колонок, по умолчанию разделитель — пробел")
	flag.Parse()

}
func main() {
	test := sortutil.SortFile()
	fmt.Print(strings.Join(test, "\r\n"))
}
