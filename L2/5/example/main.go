package main

import (
	"flag"
	"greputil"
	"log"
	"os"
)

func init() {
	if len(os.Args) < 3 {
		os.Exit(1)
	}
	greputil.FileNameGrepOpt = os.Args[len(os.Args)-1]
	greputil.ExpressionOpt = os.Args[len(os.Args)-2]

	flag.IntVar(&greputil.PrintAfterLimitOpt, "A", 0, "печатать +N строк после совпадения")
	flag.IntVar(&greputil.PrintBeforeLimitOpt, "B", 0, "печатать +N строк до совпадения")
	flag.IntVar(&greputil.PrintContextOpt, "C", 0, "(A+B) печатать ±N строк вокруг совпадения")
	flag.BoolVar(&greputil.PrintCountOpt, "c", false, "количество строк")
	//i flag.BoolVar(&greputil.IgnoreCaseOpt, "i", false, "игнорировать регистр")
	flag.BoolVar(&greputil.FilterInvertOpt, "v", false, "вместо совпадения, исключать")
	//F  flag.BoolVar(&greputil.FixedOpt, "F", false, "точное совпадение со строкой, не паттерн")
	flag.BoolVar(&greputil.PrintLineNumOpt, "n", false, "напечатать номер строки")
	flag.BoolVar(&greputil.DebugOpt, "d", false, "debug")

	flag.Parse()

}
func main() {
	err := greputil.GrepFile()
	if err != nil {
		log.Fatal(err)
	}
}
