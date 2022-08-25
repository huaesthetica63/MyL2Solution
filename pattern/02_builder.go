package pattern

import (
	"fmt"
)

/*
	Реализовать паттерн «строитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Builder_pattern
*/

/*
Паттерн "строитель" предназначен для упрощения задачи построения сложного объекта,
отделив его построение от представления. Само построение основано на использовании
простых объектов и разделения процесса на последовательные этапы.
Паттерн использует интерфейс builder,
с описанием этапов строительства, но без конкретной реализации, а также
классы concreteBuilder, которые уже реализуют конкретный способ создания сложного объекта

"Строитель" применяют, когда в системе необходимо обеспечить механизм разных вариаций схожих объектов,
сам процесс создания объектов при этом не зависит от того, из каких частей объект состоит и как они взаимосвязаны.
При этом создаваемые объекты сложные, то есть, состоят из множества компонентов и их создание
нельзя осуществить в одно действие.

Пример: нам необходимо разработать некую систему, где есть несколько групп пользователей с разными правами:
обычный пользователь (например, это потребитель, заказчик), работник (какой-либо специалист, выполняющий заказы),
администратор системы и тд. в целом, эти объекты можно назвать однотипными и похожими между собой,
но мелкие отличия все же есть - те же разные права доступа к системе, паттерн может упростить нам создание объектов-пользователей.
Другой пример: наша программа обрабатывает документы разного формата (rtf, docx, txt, листинг программ .go, .cpp и тд),
каждый из документов как объект в коде программы представляет собой сложный объект, при этом эти объекты во многом схожи,
но имеют небольшие принципиальные отличия (то же расширение). Мы здесь можем воспользоваться "строителем", как и в
первом примере

Плюсы: позволяет инкапсулировать код, делает его читабельнее и компактнее - мы можем пользоваться готовыми методами
по созданию нужных нам объектов, вместо того, чтобы разбираться в устройстве той или иной структуры, нагромождать код
большим числом строк (особенно, когда нам нужно много схожих объектов). При этом контроль над деталями так же
предоставлен - паттерн дает нам интерфейс с поэтапным описанием, как делается объект, в любой момент
мы можем изменить детали каждого из этапов под нужную нам задачу, определив новый concreteBuilder

Минусом можно назвать жесткую связь concreteBuilder с разновидностью объекта,
под каждую отдельную вариацию приходится писать новый concreteBuilder, то же самое нужно делать
при малейших изменениях в описании объекта
*/

//в качестве примера можно предемонстрировать, как "строитель" помогает упростить задачу
//создания различных объектов типа "пицца"
type Cheese struct { //сыр
	name string //наименование сыра
}
type Dough struct { //тесто
	diametr int //диаметр в см
}
type Meat struct { //колбаса
	name string //название колбасы
}
type Pizza struct { //сама пицца
	cheese  Cheese
	testo   Dough
	kolbasa Meat
}

func (p Pizza) PrintPizza() {
	fmt.Printf("Пицца: \nСыр: %s\nДиаметр в см: %d\nНачинка: %s\n", p.cheese.name, p.testo.diametr, p.kolbasa.name)
}

type BuildPizza interface { //интерфейс с описанием создания объекта
	MakeCheese() Cheese
	MakeDough() Dough
	MakeMeat() Meat
	MakePizza() Pizza // финальный этап с получением конечного сложного объекта
}
type BuildPepperoni struct { //concretebuilder для вариации объекта - пиццы Пепперони

}

func (bp BuildPepperoni) MakeCheese() Cheese {
	var res Cheese
	res.name = "Моцарелла"
	return res
}
func (bp BuildPepperoni) MakeDough() Dough {
	var res Dough
	res.diametr = 33
	return res
}
func (bp BuildPepperoni) MakeMeat() Meat {
	var res Meat
	res.name = "Салями"
	return res
}
func (bp BuildPepperoni) MakePizza() Pizza {
	var res Pizza
	res.testo = bp.MakeDough()
	res.cheese = bp.MakeCheese()
	res.kolbasa = bp.MakeMeat()
	return res
}

type BuildPizzaAndDot struct { //пицца "Пицца и точка"

}

func (bp BuildPizzaAndDot) MakeCheese() Cheese {
	var res Cheese
	res.name = "Российский"
	return res
}
func (bp BuildPizzaAndDot) MakeDough() Dough {
	var res Dough
	res.diametr = 25
	return res
}
func (bp BuildPizzaAndDot) MakeMeat() Meat {
	var res Meat
	res.name = "Ананасы"
	return res
}
func (bp BuildPizzaAndDot) MakePizza() Pizza {
	var res Pizza
	res.testo = bp.MakeDough()
	res.cheese = bp.MakeCheese()
	res.kolbasa = bp.MakeMeat()
	return res
}
func main() {
	var builder1 BuildPizzaAndDot
	pizza1 := builder1.MakePizza()
	pizza1.PrintPizza()
	var builder2 BuildPepperoni
	pizza2 := builder2.MakePizza()
	pizza2.PrintPizza()
}
