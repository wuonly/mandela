package main

import (
	"fmt"
	"math"
	"math/big"
	"math/cmplx"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"time"
)

// Outside a function, every construct begins with a keyword (var, func, and so on) and the := construct is not available
// The var statement declares a list of variables; as in function argument lists, the type is last
var x, y, z int
var c, python, java bool

// A var declaration can include initializers, one per variable
var x1, y1, z1 int = 1, 2, 3

// If an initializer is present, the type can be omitted; the variable will take the type of the initializer
var c1, python1, java1 = true, false, "no!"

// basic types
// bool
// string
// int int8 int16 int32 int64
// uint uint8 uint16 uint32 uint64 uintptr
// byte (alias for uint8)
// rune (alias for int32, represents a Unicode code point)
// float32 float64
// complex64 complex128
var (
	ToBe    bool       = false
	MaxInt  uint64     = 1<<64 - 1
	complex complex128 = cmplx.Sqrt(-5 + 12i)
)

// Constants are declared like variables, but with the const keyword
const Pi = 3.14

// Constants can be character, string, boolean, or numeric values
const World = "世界"

// Numeric Constants
const (
	Big   = 1 << 100
	Small = Big >> 99 // 2
)

type Vertex struct {
	X int
	Y int
}

type Vertex2 struct {
	Lat, Long float64
}

var m map[string]Vertex2

// Map literals are like struct literals, but the keys are required
var m2 = map[string]Vertex2{
	"gerryyang": Vertex2{
		100, 200,
	},
	"wcdj": Vertex2{
		-300, 500,
	},
}

// If the top-level type is just a type name, you can omit it from the elements of the literal
var m3 = map[string]Vertex2{
	"math":     {20, 40},
	"computer": {30, 50},
}

type Vertex3 struct {
	X, Y float64
}

type MyFloat float64

type Abser interface {
	Abs() float64
}

////////////////////////////////////////////////////////

func main() {
	fmt.Println("Hello Golang, I'm gerryyang")
	fmt.Println("The time is", time.Now())

	// To see a different number, seed the number generator; see rand.Seed
	fmt.Println("My favorite number is", rand.Intn(7))
	fmt.Printf("Now you have %g problesms\n", math.Nextafter(2, 3))
	// In Go, a name is exported if it begins with a capital letter
	fmt.Println(math.Pi)

	// Notice that the type comes after the variable name
	fmt.Println(add(42, 13))
	fmt.Println(add2(42, 13))

	// multiple results
	a, b := swap("gerry", "yang")
	fmt.Println(a, b)

	// named return
	fmt.Println(split(17))
	fmt.Println(split2(17))

	// var used
	fmt.Println(x, y, z, c, python, java)
	fmt.Println(x1, y1, z1, c1, python1, java1)

	// Inside a function, the := short assignment statement can be used in place of a var declaration with implicit type
	var x2, y2, z2 int = 1, 2, 3
	c2, python2, java2 := true, false, "yes!"
	fmt.Println(x2, y2, z2, c2, python2, java2)

	// basic types
	const f = "%T(%v)\n"
	fmt.Printf(f, ToBe, ToBe)
	fmt.Printf(f, MaxInt, MaxInt)
	fmt.Printf(f, complex, complex)

	// Constants cannot be declared using the := syntax
	const World2 = "和平"
	const Truth = true
	fmt.Println(Pi)
	fmt.Println("你好", World)
	fmt.Println("世界", World2)
	fmt.Println("Go rules?", Truth)

	// Numeric Constants
	fmt.Println(needInt(Small))
	////fmt.Println(needInt(Big))// error, constant 1267650600228229401496703205376 overflows int
	////fmt.Println(needInt64(Big)) // error, same as above
	////fmt.Println(needBigInt(big.NewInt(Big)))// error, 使用big.Int貌似入参最大类型只支持int64
	fmt.Println(needFloat(Small))
	fmt.Println(needFloat(Big))

	// Go has only one looping construct, the for loop
	// The basic for loop looks as it does in C or Java, except that the ( ) are gone (they are not even optional) and the { } are required
	sum := 0
	for i := 0; i < 10; i++ {
		sum += i
	}
	fmt.Println(sum)

	// As in C or Java, you can leave the pre and post statements empty
	// At that point you can drop the semicolons: C's while is spelled for in Go
	sum1 := 1
	for sum1 < 1000 {
		sum1 += sum1
	}
	fmt.Println(sum1)

	// If you omit the loop condition it loops forever, so an infinite loop is compactly expressed
	ivar := 1
	for {
		if ivar++; ivar > 1000 {
			fmt.Println("leave out an infinite loop")
			break
		}
	}

	// The if statement looks as it does in C or Java, except that the ( ) are gone and the { } are required
	fmt.Println(sqrt(2), sqrt(-4))

	// Like for, the if statement can start with a short statement to execute before the condition
	fmt.Println(pow(3, 2, 10), pow(3, 3, 20))

	// If and else
	fmt.Println(pow2(3, 2, 10), pow2(3, 3, 20))

	////////////////////////////////////////////////////////////

	// A struct is a collection of fields
	fmt.Println(Vertex{1, 2})

	// Struct fields are accessed using a dot
	v := Vertex{1, 2}
	v.X = 4
	fmt.Println(v)

	// Go has pointers, but no pointer arithmetic
	// Struct fields can be accessed through a struct pointer. The indirection through the pointer is transparent
	p := Vertex{1, 2}
	q := &p
	q.X = 1e9
	fmt.Println(p)

	// struct literals
	// A struct literal denotes a newly allocated struct value by listing the values of its fields
	p = Vertex{1, 2}  // has type Vertex
	q = &Vertex{1, 2} // has type * Vertex
	r := Vertex{X: 1} // Y:0 is implicit
	s := Vertex{}     // X:0 and Y:0
	fmt.Println(p, q, r, s)

	// The expression new(T) allocates a zeroed T value and returns a pointer to it
	// var t *T = new(T) or t := new(T)
	pv := new(Vertex)
	fmt.Println(pv)
	pv.X, pv.Y = 11, 9
	fmt.Println(pv)

	// A slice points to an array of values and also includes a length
	// []T is a slice with elements of type T
	as := []int{2, 3, 5, 7, 11, 13}
	fmt.Println("as ==", as)
	for i := 0; i < len(as); i++ {
		fmt.Printf("as[%d] == %d\n", i, as[i])
	}

	// Slices can be re-sliced, creating a new slice value that points to the same array
	// The expression: s[lo:hi]
	// evaluates to a slice of the elements from lo through hi-1, inclusive
	fmt.Println("as[1:4] ==", as[1:4])
	// missing low index implies 0
	fmt.Println("as[:3] ==", as[:3])
	// missing high index implies len(s)
	fmt.Println("as[4:] ==", as[4:])

	// Slices are created with the make function. It works by allocating a zeroed array and returning a slice that refers to that array
	// a := make([]int, 5), note, len(a) = 5
	// To specify a capacity, pass a third argument to make
	// b := make([]int, 0 , 5), note, len(b) = 0, cap(b) = 5
	// b = b[:cap(b)], note, len(b) = 5, cap(b) = 5
	// b = b[1:], note, len(b) = 4, cap(b) = 4
	s1 := make([]int, 5)
	printSlice("s1", s1)
	s2 := make([]int, 0, 5)
	printSlice("s2", s2)
	s3 := s2[:2]
	printSlice("s3", s3)
	s4 := s3[2:5]
	printSlice("s4", s4)

	// The zero value of a slice is nil
	// A nil slice has a length and capacity of 0
	var s5 []int
	fmt.Println(s5, len(s5), cap(s5))
	if s5 == nil {
		fmt.Println("slice is nil")
	}

	// The range form of the for loop iterates over a slice or map
	var s6 = []int{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024}
	for i, v := range s6 {
		fmt.Printf("2**%d = %d\n", i, v)
	}

	// If you only want the index, drop the ", value" entirely
	for i := range s6 {
		s6[i] = 1 << uint(i)
	}
	// You can skip the index or value by assigning to _
	for _, value := range s6 {
		fmt.Printf("%d\n", value)
	}

	// A map maps keys to values
	// Maps must be created with make (not new) before use; the nil map is empty and cannot be assigned to
	m = make(map[string]Vertex2)
	m["Bell Labs"] = Vertex2{
		40.68433, -74.39967,
	}
	fmt.Println(m["Bell Labs"])

	// Map literals are like struct literals, but the keys are required
	fmt.Println(m2)

	// If the top-level type is just a type name, you can omit it from the elements of the literal
	fmt.Println(m3)

	// map
	// insert, update, retrieve, delete, test
	m4 := make(map[string]int)
	m4["date"] = 20131129
	fmt.Println("The value:", m4["date"])
	m4["date"] = m4["date"] + 1
	fmt.Println("The value:", m4["date"])
	date, ok := m4["date"]
	fmt.Println("The value:", date, "Present?", ok)

	delete(m4, "date")
	fmt.Println("The value:", m4["date"])
	date2, ok := m4["date"]
	fmt.Println("The value:", date2, "Present?", ok)

	// Function values
	// Functions are values too
	hypot := func(x, y float64) float64 {
		return math.Sqrt(x*x + y*y)
	}
	fmt.Println(hypot(3, 4))

	// Function closures
	// For example, the adder function returns a closure. Each closure is bound to its own sum variable
	pos, neg := adder(), adder()
	for i := 0; i < 10; i++ {
		fmt.Println(pos(i), neg(-2*i))
	}

	// fibonacci
	fib := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(fib())
	}

	// switch
	// A case body breaks automatically, unless it ends with a fallthrough statement
	switch os := runtime.GOOS; os {
	case "darwin":
		fmt.Println("OS X")
	case "linux":
		fmt.Println("Linux")
	default:
		// freebsd, openbsd
		// plan9, windows...
		fmt.Printf("%s", os)
	}

	// Switch cases evaluate cases from top to bottom, stopping when a case succeeds
	// Note: Time in the Go playground always appears to start at 2009-11-10 23:00:00 UTC
	fmt.Println("When's Saturday?")
	today := time.Now().Weekday()
	switch time.Saturday {
	case today + 0:
		fmt.Println("Today")
	case today + 1:
		fmt.Println("Tomorrow")
	case today + 2:
		fmt.Println("In two days")
	case today + 3:
		fmt.Println("In three days")
	default:
		fmt.Println("Too far away")
	}

	// Switch without a condition is the same as switch true
	// This construct can be a clean way to write long if-then-else chains
	t_now := time.Now()
	switch {
	case t_now.Hour() < 12:
		fmt.Println("Good morning!")
	case t_now.Hour() < 17:
		fmt.Println("Good afternoon")
	default:
		fmt.Println("Good evening")
	}

	// Go does not have classes. However, you can define methods on struct types
	v3 := &Vertex3{3, 4}
	fmt.Println(v3.Abs())

	// In fact, you can define a method on any type you define in your package, not just structs
	// You cannot define a method on a type from another package, or on a basic type
	f1 := MyFloat(-math.Sqrt2)
	fmt.Println(f1.Abs())

	// Methods with pointer receivers
	// Methods can be associated with a named type or a pointer to a named type
	// We just saw two Abs methods. One on the *Vertex pointer type and the other on the MyFloat value type
	// There are two reasons to use a pointer receiver. First, to avoid copying the value on each method call (more efficient if the value type is a large struct). Second, so that the method can modify the value that its receiver points to
	v3 = &Vertex3{3, 4}
	v3.Scale(5)
	fmt.Println(v3, v3.Abs())

	// An interface type is defined by a set of methods
	// A value of interface type can hold any value that implements those methods
	var a_interface Abser
	v4 := Vertex3{3, 4}
	a_interface = f1  // a MyFloat implements Abser
	a_interface = &v4 // a *Vertex3 implements Abser
	//a_interface = v4  // a Vertex3, does Not, error!
	fmt.Println(a_interface.Abs())

	// Interfaces are satisfied implicitly
	var w Writer
	// os.Stdout implements Writer
	w = os.Stdout
	fmt.Fprintf(w, "hello, writer\n")

	// An error is anything that can describe itself as an error string. The idea is captured by the predefined, built-in interface type, error, with its single method, Error, returning a string: type error interface { Error() string }
	if err := run(); err != nil {
		fmt.Println(err)
	}

	// Web servers
	//var h Hello
	//http.ListenAndServe("localhost:4000", h)

}

/////////////////////////////////////////////

// func can be used before declare
func add(x int, y int) int {
	return x + y
}

// When two or more consecutive named function parameters share a type, you can omit the type from all but the last
func add2(x, y int) int {
	return x + y
}

// multiple results, a function can return any number of results
func swap(x, y string) (string, string) {
	return y, x
}

// In Go, functions can return multiple "result parameters", not just a single value. They can be named and act just like variables
func split(sum int) (x, y int) {
	x = sum * 4 / 9
	y = sum - x
	return y, x
}

// In Go, functions can return multiple "result parameters", not just a single value. They can be named and act just like variables
func split2(sum int) (x, y int) {
	x = sum * 4 / 9
	y = sum - x

	// If the result parameters are named, a return statement without arguments returns the current values of the results
	return
}

func needInt(x int) int       { return x*10 + 1 }
func needInt64(x int64) int64 { return x*10 + 1 }
func needBigInt(x *big.Int) (result *big.Int) {
	result = new(big.Int)
	result.Set(x)
	result.Mul(result, big.NewInt(10))
	return
}
func needFloat(x float64) float64 {
	return x * 0.1
}

func sqrt(x float64) string {
	if x < 0 {
		return sqrt(-x) + "i"
	}
	return fmt.Sprint(math.Sqrt(x))
}

// Variables declared by the statement are only in scope until the end of the if
func pow(x, n, lim float64) float64 {
	if v := math.Pow(x, n); v < lim {
		return v
	}
	return lim
}

// Variables declared inside an if short statement are also available inside any of the else blocks
func pow2(x, n, lim float64) float64 {
	if v := math.Pow(x, n); v < lim {
		return v
	} else {
		fmt.Printf("%g >= %g\n", v, lim)
	}
	// can't use v here, though
	return lim
}

func printSlice(s string, x []int) {
	fmt.Printf("%s len = %d cap = %d %v\n", s, len(x), cap(x), x)
}

// Go functions may be closures. A closure is a function value that references variables from outside its body. The function may access and assign to the referenced variables; in this sense the function is "bound" to the variables
func adder() func(int) int {
	sum := 0
	return func(x int) int {
		sum += x
		return sum
	}
}

// fibonacci is a function that returns
// a function that returns an int.
func fibonacci() func() int {
	p := 0
	q := 1
	s := 0
	return func() int {
		s = p + q
		p = q
		q = s
		return s
	}
}

// The method receiver appears in its own argument list between the func keyword and the method name
func (v *Vertex3) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (f MyFloat) Abs() float64 {
	if f < 0 {
		fmt.Println("f < 0 here")
		return float64(-f)
	}
	return float64(f)
}

// Methods with pointer receivers
func (v *Vertex3) Scale(f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

// Interfaces are satisfied implicitly
type Reader interface {
	Read(b []byte) (n int, err error)
}
type Writer interface {
	Write(b []byte) (n int, err error)
}
type ReadWriter interface {
	Reader
	Writer
}

// error control
type MyError struct {
	When time.Time
	What string
}

func (e *MyError) Error() string {
	return fmt.Sprintf("at %v, %s", e.When, e.What)
}
func run() error {
	return &MyError{
		time.Now(),
		"it didn't work",
	}
}

// Web servers
type Hello struct{}

func (h Hello) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request) {
	fmt.Fprint(w, "gerryyang")
}
