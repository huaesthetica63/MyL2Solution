Что выведет программа? Объяснить вывод программы.

```go
package main

import (
    "fmt"
)

func main() {
    a := [5]int{76, 77, 78, 79, 80}
    var b []int = a[1:4]
    fmt.Println(b)
}
```

Ответ:
```
Данная программа выведет срез длиной в 3 элемента, взятых из массива а: 77, 78, 79
то есть, в срез берется следующий диапазон:
первое число (перед :) - индекс первого элемента в срезе (включительно), второе число (после :) - индекс, до которого идет срез (НЕ ВКЛЮЧИТЕЛЬНО)

PS. Если бы нам потребовалось вывести и 4-ый элемент (80), то срез можно было бы записать так: a[1:] или так: a[1:5]
```