# wbl2_121223
# Паттерны
# 1 Паттерн «фасад». - 

## Основной принцип

Изменить АПИ пакета или библиотеки в сторону упрощения и лучшей юзабельности\
Предоставить более конкретный АПИ на основе более общего\
Служит как средство уменьшения связности компонентов

## Когда применяется
Когда необходимо использовать библиотеку, но не нужно создавать зависимости от сложного АПИ этой библиотеки

## Примеры из реального кода
https://github.com/labstack/echo/blob/master/echo.go \
любая верхнеровневая абстракция вероятно в некоторой степени фасад

# 2 Паттерн «строитель». - пораждающий

## Основной принцип
Позволяет кастомизировать "конструктор" экземпляра. Очень часто встречается. Основной признак того что перед нами билдер - WithМетоды в интерфейсе, которые возвращают тот же интерфейс, который реализуют. func (i iface1)WithSomething iface. \
Это позволяет удобно контролировать набор опций, без уточнения какие опции конкретно будут использоваться. Например можно создать объект с Валидатором, Логгером, Аутентификатором.

## Когда применяется

Когда понятен набор опций, он большой, и большое количество непредопределенных комбинаций использования этих опций

## Примеры из реального кода
https://github.com/gotd/td/blob/main/telegram/message/sender.go

# 3 Паттерн «посетитель». - поведенческий

## Основной принцип
Позволяет разделить алгоритм от структуры объекта.\
Объект реализует интерфейс, который принимает посетителя.\
Посетитель знает как работать с конкретным объектом и реализует какую то логику.

## Когда применяется
Когда есть определенный набор однотипных действий, которые необходимо производить над разородными объектами

Например есть 
```
type textPayload struct{ 
	b []byte

}
type picturePayload struct{ 
	b []byte
}

//И есть 
type processor interface{
	ProcessText(p textPayload)
	ProcessPic(p picturePayload)
}
```
очевидно, что каждая реализация payload будет пользоваться разными методами processor
```
func (t textPayload)AcceptProcessor(v processor){
	v.ProcessText(t)
}
func (p picturePayload)AcceptProcessor(v processor){
	v.ProcessPic(p)
}
```
процессоров может быть сколько угодно, от ПДФ конвертора, который умеет работать как с текстом так и с картиками, до какого нибудь ну другого конвертора или сортировщика\
очевидно что для пущей эффективности можно описать payload интерфейс, но для посетителя это не принципально

## Примеры из реального кода
https://github.com/gotd/td/blob/main/telegram/message/sender.go \
как ни странно снова этот пример, но теперь функции sendMessage / sendMedia / sendMultiMedia

# 4 Паттерн «комманда». - поведенческий

## Основной принцип
Позволяет выделить схожие метода встроенного объекта в отдельный объект с унифицированным интерфейсом\
Почти как посетитель но вывернутый наизнанку

## Когда применяется
Когда необходимо упаковать методы объекта в отдельные унифицированные команды, для их последующего использования
```
package main

import "fmt"

type button struct {
    command command
}

func (b *button) press() {
    b.command.execute()
}

type command interface {
    execute()
}

type offCommand struct {
    device device
}

func (c *offCommand) execute() {
    c.device.off()
}

type onCommand struct {
    device device
}

func (c *onCommand) execute() {
    c.device.on()
}

type device interface {
    on()
    off()
}

type tv struct {
    isRunning bool
}

func (t *tv) on() {
    t.isRunning = true
    fmt.Println("Turning tv on")
}

func (t *tv) off() {
    t.isRunning = false
    fmt.Println("Turning tv off")
}

func main() {
    tv := &tv{}
    onCommand := &onCommand{
        device: tv,
    }
    offCommand := &offCommand{
        device: tv,
    }
    onButton := &button{
        command: onCommand,
    }
    onButton.press()
    offButton := &button{
        command: offCommand,
    }
    offButton.press()
}
```
## Примеры из реального кода
Пока не найдено

# 5 Паттерн «цепочка вызовов» - поведенческий

## Основной принцип
Позволяет создавать цепочку обработчиков, где каджых входящий запрос будет проходить по цепочке с вызовом каждого обработчика

## Когда применяется
Когда можно выделить отдельные этапы обработки запроса с единым интерфейсом\
По сути создается связный список из реализаций обработчиков, с методом execute, который запускает цепочку вызовов. (траверс по связному списку)
```
package main

import "fmt"

type department interface {
    execute(*patient)
    setNext(department)
}

type reception struct {
    next department
}

func (r *reception) execute(p *patient) {
    if p.registrationDone {
        fmt.Println("Patient registration already done")
        r.next.execute(p)
        return
    }
    fmt.Println("Reception registering patient")
    p.registrationDone = true
    r.next.execute(p)
}

func (r *reception) setNext(next department) {
    r.next = next
}

type doctor struct {
    next department
}

func (d *doctor) execute(p *patient) {
    if p.doctorCheckUpDone {
        fmt.Println("Doctor checkup already done")
        d.next.execute(p)
        return
    }
    fmt.Println("Doctor checking patient")
    p.doctorCheckUpDone = true
    d.next.execute(p)
}

func (d *doctor) setNext(next department) {
    d.next = next
}

type medical struct {
    next department
}

func (m *medical) execute(p *patient) {
    if p.medicineDone {
        fmt.Println("Medicine already given to patient")
        m.next.execute(p)
        return
    }
    fmt.Println("Medical giving medicine to patient")
    p.medicineDone = true
    m.next.execute(p)
}

func (m *medical) setNext(next department) {
    m.next = next
}

type cashier struct {
    next department
}

func (c *cashier) execute(p *patient) {
    if p.paymentDone {
        fmt.Println("Payment Done")
    }
    fmt.Println("Cashier getting money from patient patient")
}

func (c *cashier) setNext(next department) {
    c.next = next
}

type patient struct {
    name              string
    registrationDone  bool
    doctorCheckUpDone bool
    medicineDone      bool
    paymentDone       bool
}

func main() {
    cashier := &cashier{}
   
    //Set next for medical department
    medical := &medical{}
    medical.setNext(cashier)
   
    //Set next for doctor department
    doctor := &doctor{}
    doctor.setNext(medical)
   
    //Set next for reception department
    reception := &reception{}
    reception.setNext(doctor)
   
    patient := &patient{name: "abc"}
    //Patient visiting
    reception.execute(patient)
}
```
## Примеры из реального кода
Пока не найдено

# 6 Паттерн «фабричный метод» - порождающий

## Основной принцип
Удобно хранит все конструкторы объекта в одном объекте, методе или функции\
Позволяет создавать объект не зная его конкретный тип, получая сразу интерфейс\
Легко расширить новыми конструкторами\
Разделяет процесс создания объекта от его использования

## Когда применяется
Когда необходимо организовать конструкторы в одном месте и возвращать интерфейс

## Примеры из реального кода 
https://github.com/gotd/td/blob/main/telegram/message/peer/resolver.go

# 7 Паттерн «стратегия» - поведенческий

## Основной принцип
Делегируем основную логику поведения объекта вложенному интерфейсу\
Так же как и в "состоянии", разница в том, что объект не реализует вложенный интерфейс, а передается через методы в интерфейс. Тем самым может реализовать другой интерфейс.\
Что-то вроде смеси адаптера и состояния\

## Когда применяется
Когда нужно определить поведение объекта несовместимым интерфейсом

## Примеры из реального кода
Пока не найдено

# 8 Паттерн «состояние» - поведенческий

## Основной принцип
Основано на конечных автоматах, поведение которых определяет состояние в котором находится конечный автомат\
За счет композиции, любая структура со вложенным интерфейсом автоматически реализует этот интерфейс, и получается State, так что пример приводить особо нет смысла

## Когда применяется
Практически интегрирован в язык за счет композиции\
Когда объект может находиться в разных состояниях и в зависимости от них изменять свое поведение

## Примеры из реального кода
