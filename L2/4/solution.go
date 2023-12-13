/*
Данный фрагмент написан на языке Golang.
Он реализует алгоритм нахождения анаграмм, т.е. слов с одинаковой последовательностью букв в заданном наборе слов.
(комментарии любезно предоставлены YandexGPT2 :)
*/
package anogram

import (
	"sort"
	"strings"
)

/*
функция, принимающая на вход массив строк и возвращающая карту с хэшами этих строк в виде ключей и отсортированными списками этих строк в виде значений.
*/
func AnogramsFromSlice(s []string) map[string]words {
	res := make(map[string]words)
	//тут всё очевидно
	for _, v := range s {
		//приводим к нижнему регистру
		v = strings.ToLower(v)
		//извлекаем по хэшу накопленную последовательность
		sl := res[hash(v)]
		//проверяем на повторы
		if !hasElement(sl, v) {
			sl = append(sl, (v))
			//кладем по бакетам - хэшам строки
			res[hash(v)] = sl
		}
	}
	for k, v := range res {
		//получили исходный срез
		w := words(v)
		//удалили его сразу
		delete(res, k)
		if len(w) < 2 {
			//если он слишком маленький продолжаем следующую итерацию
			continue
		}
		//извлекаем первый элемент для ключа
		first := w[0]
		//сортируем анограмы
		sort.Sort(w)
		//устанавливаем отсортированные по ключу
		res[first] = w
	}
	return res
}

/*
функция для проверки наличия элемента в списке.
Принимает на вход список и элемент.
Возвращает true, если элемент есть в списке и false, если нет.
*/
func hasElement[T comparable](s []T, e T) bool {
	for i := range s {
		if e == s[i] {
			return true
		}
	}
	return false
}

/*
hash - функция хэширования строки.
Принимает строку и возвращает строку, представляющую собой хэш этой строки.
*/
func hash(s string) string {
	l := letters(s)
	sort.Sort(l)
	return string(l)
}

/*
letters - структура данных для работы со строками.
Позволяет выполнять операции над строками, такие как сортировка, сравнение и т.д.
*/
type letters []rune

func (l letters) Len() int {
	return len(l)
}

func (l letters) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
func (l letters) Less(i, j int) bool {
	return l[i] < l[j]
}

/*
words - структура данных, представляющая собой список строк.
Используется для хранения отсортированных списков слов.
*/
type words []string

func (w words) Len() int {
	return len(w)
}

func (w words) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
}
func (w words) Less(i, j int) bool {
	return w[i] < w[j]
}
