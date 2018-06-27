package vm_test

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/mattn/anko/vm"
)

func ExampleInterrupt() {
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	waitChan := make(chan struct{}, 1)

	env := vm.NewEnv()
	sleepMillisecond := func() { time.Sleep(time.Millisecond) }

	err := env.Define("println", fmt.Println)
	if err != nil {
		log.Fatalf("define error: %v\n", err)
	}
	err = env.Define("sleep", sleepMillisecond)
	if err != nil {
		log.Fatalf("define error: %v\n", err)
	}

	script := `
# sleep for 10 seconds
for i = 0; i < 10000; i++ {
	sleep()
}
# Should interrupt before printing the next line
println("this line should not be printed")
`

	go func() {
		close(waitChan)
		v, err := env.Execute(script)
		fmt.Println(v, err)
		waitGroup.Done()
	}()

	<-waitChan
	vm.Interrupt(env)

	waitGroup.Wait()

	// output: <nil> execution interrupted
}

func ExampleEnv_Define() {
	env := vm.NewEnv()
	env.SetName("myName")

	err := env.Define("println", fmt.Println)
	if err != nil {
		log.Fatalf("define error: %v\n", err)
	}

	err = env.Define("a", true)
	if err != nil {
		log.Fatalf("define error: %v\n", err)
	}
	err = env.Define("b", int64(1))
	if err != nil {
		log.Fatalf("define error: %v\n", err)
	}
	err = env.Define("c", float64(1.1))
	if err != nil {
		log.Fatalf("define error: %v\n", err)
	}
	err = env.Define("d", "d")
	if err != nil {
		log.Fatalf("define error: %v\n", err)
	}
	err = env.Define("e", []interface{}{true, int64(1), float64(1.1), "d"})
	if err != nil {
		log.Fatalf("define error: %v\n", err)
	}
	err = env.Define("f", map[string]interface{}{"a": true})
	if err != nil {
		log.Fatalf("define error: %v\n", err)
	}

	script := `
println(a)
println(b)
println(c)
println(d)
println(e)
println(f)
`

	_, err = env.Execute(script)
	if err != nil {
		log.Fatalf("execute error: %v\n", err)
	}

	// output:
	// true
	// 1
	// 1.1
	// d
	// [true 1 1.1 d]
	// map[a:true]
}

func ExampleEnv_Dump() {
	env := vm.NewEnv()
	env.SetName("myName")

	err := env.Define("a", "a")
	if err != nil {
		log.Fatalf("define error: %v\n", err)
	}

	_, err = env.Get("a")
	if err != nil {
		log.Fatalf("get error: %v\n", err)
	}

	env.Dump()

	// output:
	// Name: myName
	// Has parent: false
	// a = "a"
}

func Example_vmHelloWorld() {
	env := vm.NewEnv()

	err := env.Define("println", fmt.Println)
	if err != nil {
		log.Fatalf("define error: %v\n", err)
	}

	script := `
println("Hello World :)")
`

	_, err = env.Execute(script)
	if err != nil {
		log.Fatalf("execute error: %v\n", err)
	}

	// output: Hello World :)
}

func Example_vmQuickStart() {
	env := vm.NewEnv()

	err := env.Define("println", fmt.Println)
	if err != nil {
		log.Fatalf("define error: %v\n", err)
	}

	script := `
// declare variables
x = 1
y = x + 1

// print using outside the script defined println function
println(x + y) // 3

// if else statement
if x < 1 || y < 1 {
	println(x)
} else if x < 1 && y < 1 {
	println(y)
} else {
	println(x + y)
}

// array
a = [1, 2, 3]
println(a) // [1 2 3]
println(a[0]) // 1

// map
a = {"x": 1}
println(a) // map[x:1]
a.b = 2
a["c"] = 3
println(a["b"]) // 2
println(a.c) // 3

// function
func a (x) {
	println(x + 1)
}
a(3) // 4
`

	_, err = env.Execute(script)
	if err != nil {
		log.Fatalf("execute error: %v\n", err)
	}

	// output:
	// 3
	// 3
	// [1 2 3]
	// 1
	// map[x:1]
	// 2
	// 3
	// 4
}
