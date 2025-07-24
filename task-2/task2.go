package main

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

/*指针开始*/
// 编写一个Go程序，定义一个函数，该函数接收一个整数指针作为参数，在函数内部将该指针指向的值增加10，然后在主函数中调用该函数并输出修改后的值
func addTen(num *int) {
	*num += 10
}

// 实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
func multiplyByTwo(slice *[]int) {
	for i := range *slice {
		(*slice)[i] *= 2
	}
}

/*指针结束*/
/*goroutine开始*/
// 编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
func printNumbers() {
	var wg sync.WaitGroup
	wg.Add(2)
	// 打印奇数的协程
	go func() {
		defer wg.Done()
		for i := 1; i <= 10; i += 2 {
			fmt.Printf("奇数: %d\n", i)
			time.Sleep(100 * time.Millisecond)
		}
	}()
	// 打印偶数的协程
	go func() {
		defer wg.Done()
		for i := 2; i <= 10; i += 2 {
			fmt.Printf("偶数: %d\n", i)
			time.Sleep(100 * time.Millisecond)
		}
	}()
	wg.Wait()
}

// 题设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
type Task func() string

func taskScheduler(tasks []Task) {
	var wg sync.WaitGroup
	wg.Add(len(tasks))
	for i, task := range tasks {
		go func(id int, t Task) {
			defer wg.Done()
			start := time.Now()
			result := t()
			elapsed := time.Since(start)
			fmt.Printf("任务 %d 完成，结果: %s，耗时: %v\n", id, result, elapsed)
		}(i, task)
	}
	wg.Wait()
}

/*goroutine结束*/

/*面向对象开始*/
// 题目1：Shape接口及其实现
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Rectangle结构体
type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// Circle结构体
type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// 题目2：组合示例
type Person struct {
	Name string
	Age  int
}

type Employee struct {
	Person     // 组合Person结构体
	EmployeeID string
}

func (e Employee) PrintInfo() {
	fmt.Printf("员工ID: %s, 姓名: %s, 年龄: %d\n", e.EmployeeID, e.Name, e.Age)
}

/*面向对象结束*/

/*channel开始*/

// 题目1：使用通道实现协程间通信
func channelCommunication() {
	ch := make(chan int)
	var wg sync.WaitGroup
	wg.Add(2)

	// 生产者协程
	go func() {
		defer wg.Done()
		defer close(ch)
		for i := 1; i <= 10; i++ {
			ch <- i
			//发送数据打印
			fmt.Printf("发送数字: %d\n", i)
		}
	}()

	// 消费者协程
	go func() {
		defer wg.Done()
		for num := range ch {
			fmt.Printf("接收到数字: %d\n", num)
		}
	}()

	wg.Wait()
}

// 题目2：带缓冲的通道
func bufferedChannel() {
	ch := make(chan int, 10) // 缓冲大小为10
	var wg sync.WaitGroup
	wg.Add(2)

	// 生产者协程
	go func() {
		defer wg.Done()
		defer close(ch)
		for i := 1; i <= 100; i++ {
			ch <- i
			//发送数据打印
			fmt.Printf("发送数字: %d\n", i)
		}
	}()

	// 消费者协程
	go func() {
		defer wg.Done()
		count := 0
		for num := range ch {
			count++
			fmt.Printf("接收到第 %d 个数字: %d\n", count, num)
		}
	}()

	wg.Wait()
}

/*channel结束*/

/*锁机制开始*/

// 题目1：使用互斥锁保护计数器
func mutexCounter() {
	var counter int
	var mutex sync.Mutex
	var wg sync.WaitGroup

	// 启动10个协程
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				mutex.Lock()
				counter++
				mutex.Unlock()
			}
		}()
	}

	wg.Wait()
	fmt.Printf("使用互斥锁的计数器最终值: %d\n", counter)
}

// 题目2：使用原子操作实现无锁计数器
func atomicCounter() {
	var counter int64
	var wg sync.WaitGroup

	// 启动10个协程
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				atomic.AddInt64(&counter, 1)
			}
		}()
	}

	wg.Wait()
	fmt.Printf("使用原子操作的计数器最终值: %d\n", counter)
}

/*锁机制结束*/

func main() {
	// 将指针指向的值增加10
	num := 5
	fmt.Printf("原始值: %d\n", num)
	addTen(&num)
	fmt.Printf("增加10后: %d\n", num)

	// 将切片中的每个元素乘以2
	slice := []int{1, 2, 3, 4, 5}
	fmt.Printf("原始切片: %v\n", slice)
	multiplyByTwo(&slice)
	fmt.Printf("每个元素乘以2后: %v\n", slice)

	// 打印奇数和偶数
	fmt.Println("=== 测试打印奇数和偶数 ===")
	printNumbers()

	// 任务调度器
	fmt.Println("\n=== 测试任务调度器 ===")
	tasks := []Task{
		func() string {
			time.Sleep(100 * time.Millisecond)
			return "任务1完成"
		},
		func() string {
			time.Sleep(200 * time.Millisecond)
			return "任务2完成"
		},
		func() string {
			time.Sleep(150 * time.Millisecond)
			return "任务3完成"
		},
	}
	taskScheduler(tasks)

	// Shape接口及其实现
	fmt.Println("=== 测试Shape接口 ===")
	rect := Rectangle{Width: 5, Height: 3}
	circle := Circle{Radius: 2.5}

	shapes := []Shape{rect, circle}
	for i, shape := range shapes {
		if i == 0 {
			fmt.Println("矩形:")
		} else {
			fmt.Println("圆形:")
		}
		fmt.Printf("  面积: %.2f\n", shape.Area())
		fmt.Printf("  周长: %.2f\n", shape.Perimeter())
	}

	// 组合示例
	fmt.Println("\n=== 测试组合 ===")
	emp := Employee{
		Person:     Person{Name: "张三", Age: 30},
		EmployeeID: "EMP001",
	}
	emp.PrintInfo()

	// channel
	fmt.Println("=== 测试通道通信 ===")
	channelCommunication()

	// channel带缓冲
	fmt.Println("\n=== 测试带缓冲的通道 ===")
	bufferedChannel()

	// 测试题目1
	fmt.Println("=== 测试互斥锁 ===")
	mutexCounter()

	// 测试题目2
	fmt.Println("\n=== 测试原子操作 ===")
	atomicCounter()

}
