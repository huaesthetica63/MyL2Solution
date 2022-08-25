package pattern

/*
	Реализовать паттерн «посетитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Visitor_pattern
*/

/*
Паттерн "посетитель" применяется, когда у нас есть объект, и к нему нужно добавить новую функциональность, не
меняя при этом структуру самого объекта, чтобы его не загромождать. Другой причиной использования паттерна
может являться постоянное изменение функционала объектов. Например, у нас есть какая-то информация о заказе
в какой-то момент нам может понадобиться функция сериализации в json, чуть позже в html, еще чуть позже в xml,
потом может потребоваться формирование отчетов в xlsx или pdf, и так, постоянно мы можем добавлять все новые и
новые методы к структуре с описанием заказа, что сделает код, ее описывающий, очень громоздким, вместо этого
можно создать отдельно visitor, который будет предназначен для всех этих методов, и мы будем дополнять только его,
не касаясь структуры заказа. Также, можно предположить, что паттерн "посетитель" будет удобным в ситуации,
когда у нас есть некий код с описанием объектов, который часто используется в разнообразных проектах, в одном
из таких проектов от объекта нужна функциональность, которая больше нигде кроме этого проекта применяться не будет -
слишком узконаправленная задача. Например, у нас есть библиотека для построения графиков математических функций,
в одном из проектов зачем-то нужно вывести в лог-файл значения функции на интервале [-100, 100]. Поскольку этот
метод нужен исключительно в этом проекте и вряд ли будет полезен в самой библиотеке, мы и воспользуемся помощью паттерна

Плюсы: упрощает добавление новых операций и объединение родственных методов, облегчает понимание и использование кода,
потому что сами структуры, описывающие объекты, не нагромождены множеством методов.
К недостаткам можно отнести то, что добавление новых объектов, которым тоже нужен visitor, затрудняет использование паттерна,
придется прописать достаточно много кода, если вдруг иерархия объектов поменяется. Поэтому можно сделать вывод,
что "посетитель" уместнее в ситуации, когда система часто меняет или добавляет функционал, но не сами объекты
*/
import (
	"fmt"
	"math"
)

//реализуем паттерн на примере геометрических фигур
type Rectangle struct { //прямоугольник
	a, b float64 //различные стороны прямоугольника
}
type Triangle struct { //треугольник
	a, b, c float64
}
type Circle struct { //окружность
	r float64 //радиус
}
type Visitor interface { //интерфейс посетителя
	VisitRectangle(r Rectangle)
	VisitTriangle(t Triangle)
	VisitCircle(c Circle)
}
type Figure interface { //этот интерфейс должны реализовывать все наши фигуры, чтобы подходить под паттерн
	Accept(Visitor) //accept метод принимает "визитера"
}

func (r Rectangle) Accept(v Visitor) {
	v.VisitRectangle(r)
}
func (t Triangle) Accept(v Visitor) {
	v.VisitTriangle(t)
}
func (c Circle) Accept(v Visitor) {
	v.VisitCircle(c)
}

//теперь весь новый функционал можно добавлять в отдельные структуры отдельно для каждого объекта
type CalculPerimeter struct { //новый функционал - посчитаем периметр (для окружности длину окружности)
	perimeter float64
}

func (c *CalculPerimeter) VisitRectangle(r Rectangle) {
	c.perimeter = r.a*2 + r.b*2
}
func (c *CalculPerimeter) VisitTriangle(t Triangle) {
	c.perimeter = t.a + t.b + t.c
}
func (c *CalculPerimeter) VisitCircle(cir Circle) {
	c.perimeter = 2 * cir.r * math.Pi
}

type CalculSquare struct { //новый функционал - посчитаем площадь
	square float64
}

func (c *CalculSquare) VisitRectangle(r Rectangle) {
	c.square = r.a * r.b
}
func (c *CalculSquare) VisitTriangle(t Triangle) {
	//для площади треугольника применим формулу Герона
	p := (t.a + t.b + t.c) / 2 //полупериметр
	c.square = math.Sqrt(p * (p - t.a) * (p - t.b) * (p - t.c))
}
func (c *CalculSquare) VisitCircle(cir Circle) {
	c.square = cir.r * cir.r * math.Pi
}

func main() {
	rect := Rectangle{a: 10, b: 20}
	triangle := Triangle{a: 30, b: 25, c: 10}
	circ := Circle{r: 10}
	var cp CalculPerimeter
	var cs CalculSquare
	fmt.Println("Прямоугольник: ", rect)
	rect.Accept(&cp)
	rect.Accept(&cs)
	fmt.Printf("Периметр: %f\n", cp.perimeter)
	fmt.Printf("Площадь: %f\n", cs.square)
	fmt.Println("Треугольник: ", triangle)
	triangle.Accept(&cp)
	triangle.Accept(&cs)
	fmt.Printf("Периметр: %f\n", cp.perimeter)
	fmt.Printf("Площадь: %f\n", cs.square)
	fmt.Println("Окружность: ", circ)
	circ.Accept(&cp)
	circ.Accept(&cs)
	fmt.Printf("Периметр: %f\n", cp.perimeter)
	fmt.Printf("Площадь: %f\n", cs.square)

}
