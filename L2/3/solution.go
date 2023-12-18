package sortutil

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// функция сортирует файл с конфигурацией из экспортируемых переменных
func SortFile() []string {
	//чтение файла
	b, err := os.ReadFile(FileNameSortOpt)
	if err != nil {
		panic(err)
	}
	if DebugSortOpt {
		fmt.Println("debug:", "file size", len(b))
	}
	//разбивание на строки
	lines := sortBy(strings.Split(string(b), "\n"))
	if DebugSortOpt {
		fmt.Println("debug:", "file lines", len(lines))
	}
	//очистка от символа возврата каретки (к сожалению сейчас необходимо юзать винду)
	for i := range lines {
		lines[i] = strings.TrimSuffix(lines[i], "\r")
	}
	//фильтруем лишнее
	lines = applyFilterOpts(lines)
	if DebugSortOpt {
		fmt.Println("debug:", "file lines after filter", len(lines))
	}
	//сортировка
	sort.Sort(lines)
	if DebugSortOpt {
		fmt.Println("debug:", "file lines after sort", len(lines))
	}
	return lines
}

// сортирующийся срез строк
type sortBy []string

var (
	UniqueSortOpt          = true
	ReverseSortOpt         = false
	CaseSensetiveSortOpt   = true
	ColumnSortOpt          = 2
	ColumnSeparatorSortOpt = " "
	NumberSortOpt          = true
	FileNameSortOpt        = "test_data.txt"
	DebugSortOpt           = true
)

// длина
func (a sortBy) Len() int {
	return len(a)
}

// замена
func (a sortBy) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// сравнение
//
// тут мы применяем настройки сравнения и преобразования данных
func (a sortBy) Less(i, j int) bool {
	return applyCompareOpts(applyDataOpts(a[i]), applyDataOpts(a[j]))
}

// функция фильтр. фильтрует уникальные значения если есть флаг
func applyFilterOpts(s []string) []string {
	result := make([]string, 0, len(s))
	for _, v := range s {
		if contains(result, v) && UniqueSortOpt {
			if DebugSortOpt {
				fmt.Println("debug:", "filter:", "filtered line", v)
			}
			continue
		}
		result = append(result, v)
	}
	return result
}

// функция сравнивает в зависимости от флагов. сравнивает как число или как строку
func applyCompareOpts(si, sj string) bool {
	if NumberSortOpt {
		ni, err := strconv.ParseInt(si, 10, 64)
		if err != nil {
			panic(err)
		}
		nj, err := strconv.ParseInt(sj, 10, 64)
		if err != nil {
			panic(err)
		}
		return compareWithDirection(ni, nj)
	}
	return compareWithDirection(si, sj)
}

// функция применяет реверс если есть флаг
func compareWithDirection[T int64 | string](i, j T) bool {
	if !ReverseSortOpt {
		return i < j
	} else {
		return i > j
	}
}

// функция преобразования данных для сравнения. нужна чтобы вычленить столбец из строки и применить настроку нечувствительности к регистру
func applyDataOpts(s string) string {
	words := strings.Split(s, ColumnSeparatorSortOpt)
	if ColumnSortOpt >= 0 && len(words) > ColumnSortOpt {
		s = words[ColumnSortOpt]
	}
	if !CaseSensetiveSortOpt {
		s = strings.ToLower(s)
	}
	return s
}

// функция помошник для определения содержится ли элемент в срезе
func contains[T comparable](arr []T, elem T) bool {
	for _, v := range arr {
		if v == elem {
			return true
		}
	}
	return false
}
