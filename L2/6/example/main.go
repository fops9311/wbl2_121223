package main

import (
	"bufio"
	"cututil"
	"flag"
	"fmt"
	"log"
	"os"
)

func init() {
	if len(os.Args) < 3 {
		os.Exit(1)
	}
	cututil.CutFileName = os.Args[len(os.Args)-1]

	flag.IntVar(&cututil.CutColumn, "f", 0, "выбрать поля (колонки)")
	flag.BoolVar(&cututil.CutSeparatedOpts, "s", false, "только строки с разделителем")
	flag.StringVar(&cututil.CutSepOpt, "d", "\t", "использовать другой разделитель")

	flag.Parse()

}
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	r, err := cututil.CutFile(scanner)
	if err != nil {
		log.Fatal(err)
		return
	}
	outScanner := bufio.NewScanner(r)
	for outScanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}
	if outScanner.Err() != nil {
		log.Fatal(err)
		return
	}
}
