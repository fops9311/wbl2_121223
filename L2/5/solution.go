package greputil

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	FileNameGrepOpt               = "test_data.txt"
	DebugOpt                      = true
	PrintLineNumOpt               = true
	PrintBeforeLimitOpt           = 0
	PrintAfterLimitOpt            = 0
	PrintContextOpt               = 1
	PrintCountOpt                 = true
	FilterInvertOpt               = false
	MatchModeOpt        FilterOpt = FixedMatchFilterOpt
	ExpressionOpt                 = "laptop"
)

func PrintApplyContext(v int) int {
	if PrintContextOpt > v {
		return PrintContextOpt
	}
	return v
}

// функция фильтрует файл с конфигурацией из экспортируемых переменных
func GrepFile() error {
	//чтение файла
	file, err := os.Open(FileNameGrepOpt)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	in := make(chan string)
	chanscanner := &ChanScanner{
		Scanner:       scanner,
		index:         0,
		printIndexOpt: PrintLineNumOpt,
	}
	go chanscanner.ScanToChan(in)
	sf, err := NewSearchFilter(ExpressionOpt, MatchModeOpt)
	if err != nil {
		return err
	}
	PrintBeforeLimitOpt = PrintApplyContext(PrintBeforeLimitOpt)
	PrintAfterLimitOpt = PrintApplyContext(PrintAfterLimitOpt)
	printer := &Printer{
		printBeforeLimit: PrintBeforeLimitOpt,
		printAfterLimit:  PrintAfterLimitOpt,
	}
	out, bad := Filter(in, sf.GetFilterFunc())
	var printfunc func(a ...any) (n int, err error)
	var count int
	if PrintCountOpt {
		printfunc = func(a ...any) (n int, err error) {
			return
		}
		defer func() { fmt.Println("count:", count) }()
	} else {
		printfunc = fmt.Println
	}

	count, err = printer.PrintFromChan(out, bad, printfunc)
	if err != nil {
		return err
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

type Printer struct {
	queue             []string
	printBeforeLimit  int
	printAfterLimit   int
	printAfterCounter int
}

func (p *Printer) PrintFromChan(good <-chan string, bad <-chan string, f func(a ...any) (n int, err error)) (int, error) {
	var count int
	p.queue = make([]string, 0, p.printBeforeLimit+1)
	printer := func(a ...any) error {
		_, err := f(a...)
		if err != nil {
			return err
		}
		count++
		return nil
	}
	for {
		select {
		case s, ok := <-good:
			if !ok {
				good = nil
				if bad == nil {
					return count, nil
				}
				continue
			}
			for 0 < len(p.queue) {
				var bad string
				bad, p.queue = dequeue(p.queue)
				err := printer("-", bad)
				if err != nil {
					return count, err
				}
			}
			p.printAfterCounter = p.printAfterLimit
			err := printer("+", s)
			if err != nil {
				return count, err
			}
		case s, ok := <-bad:
			if !ok {
				bad = nil
				if good == nil {
					return count, nil
				}
				continue
			}
			if p.printAfterCounter > 0 {
				p.printAfterCounter--
				err := printer("-", s)
				if err != nil {
					return count, err
				}
				continue
			}
			p.queue = enqueue(p.queue, s)
			if len(p.queue) > p.printBeforeLimit {
				_, p.queue = dequeue(p.queue)
			}
		}
	}
}

type ChanScanner struct {
	*bufio.Scanner
	index         int
	printIndexOpt bool
}

func (cs *ChanScanner) ScanToChan(in chan<- string) {
	if cs.printIndexOpt {
		for cs.Scan() {
			in <- fmt.Sprintf("%d: %s", cs.index, cs.Text())
			cs.index++
		}
	} else {
		for cs.Scan() {
			in <- cs.Text()
		}
	}
	close(in)
}

type SearchFilter struct {
	*regexp.Regexp
	Expr string
	Opt  FilterOpt
}
type FilterOpt int

const (
	FixedMatchFilterOpt FilterOpt = iota
	RegExpMatchFilterOpt
)

func NewSearchFilter(expr string, opt FilterOpt) (SearchFilter, error) {
	if opt == RegExpMatchFilterOpt {
		re, err := regexp.Compile(expr)
		if err != nil {
			return SearchFilter{}, err
		}
		return SearchFilter{
			Regexp: re,
			Opt:    opt,
		}, nil
	}
	return SearchFilter{
		Expr: expr,
		Opt:  opt,
	}, nil
}
func (sf SearchFilter) GetFilterFunc() func(s string) bool {
	switch sf.Opt {
	case FixedMatchFilterOpt:
		return sf.FixedMatchFilter
	case RegExpMatchFilterOpt:
		return sf.RegExpMatchFilter
	default:
		return func(s string) bool { return true }
	}
}

func (sf SearchFilter) FixedMatchFilter(s string) bool {
	return strings.Contains(s, sf.Expr)
}
func (sf SearchFilter) RegExpMatchFilter(s string) bool {
	return sf.MatchString(s)
}

//добавляем в связных список
//очистка списка вызывается раз в 2 длины списка
//добавить функцию дамп последних N элементов списка
//список обрабатывается конечным автоматом
//если в режиме поиска находим совпадение - делаем дамп и начинаем считать стрик
//если во время стрика находим совпадение - переинициализация счетчика
//если нет, уменьшаем счетчик
//если счетчик дошел до нуля переходим обратно в режим поиска

func Filter(in <-chan string, f func(s string) bool) (out chan string, bad chan string) {
	out = make(chan string)
	bad = make(chan string)
	filterfunc := func() func(s string) bool {
		if FilterInvertOpt {
			return func(s string) bool {
				return !f(s)
			}
		} else {
			return func(s string) bool {
				return f(s)
			}
		}
	}()
	go func() {
		for s := range in {
			if filterfunc(s) {
				out <- s
			} else {
				bad <- s
			}
		}
		close(out)
		close(bad)
	}()
	return
}

func enqueue[T any](queue []T, element T) []T {
	queue = append(queue, element) // Simply append to enqueue.
	return queue
}

func dequeue[T any](queue []T) (T, []T) {
	element := queue[0] // The first element is the one to be dequeued.
	if len(queue) == 1 {
		var tmp = []T{}
		return element, tmp

	}
	return element, queue[1:] // Slice off the element once it is dequeued.
}
