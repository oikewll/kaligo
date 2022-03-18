package kaligo

import (
    "net/http"
    "testing"

    "github.com/owner888/kaligo/routes"
    "github.com/stretchr/testify/assert"
)

var _ routes.Router = &KaliGo{}
var _ http.Handler = &KaliGo{}

func TestNew(t *testing.T) {
    kali := New()
    assert.NotNil(t, kali)
}

// 1、其他编程语言不常看到的：把map的value设置成函数，可以实现工厂模式
// 2、golang不具备set，但是可以用map来实现
//func TestMapWithFuncValue(t *testing.T) {
//m := map[int]func(op int) int{}
//m[1] = func(op int) int { return op }
//m[2] = func(op int) int { return op * op }
//m[3] = func(op int) int { return op * op * op }
//t.Log(m[1](2), m[2](2), m[3](2))
//}

// 实现Set
//map[type]bool
// 基本操作
// 1、添加元素
// 2、判断元素是否存在
// 3、删除元素
// 4、元素个数
//func TestMapForSet(t *testing.T) {
//mySet := map[int]bool{}
//mySet[1] = true
//n := 1
//if mySet[n] {
//t.Logf("%d is existing", n)
//} else {
//t.Logf("%d is not existing", n)
//}

//mySet[3] = true
//t.Log(len(mySet))
//delete(mySet, 1)
//n = 1
//if mySet[n] {
//t.Logf("%d is existing", n)
//} else {
//t.Logf("%d is not existing", n)
//}
//}

//func Sum(ops ...int) int {
//ret := 0
//for _, op := range ops {
//ret += op
//}
//return ret
//}

//func TestVarParam(t *testing.T) {
//t.Log(Sum(1,2,3,4))
//t.Log(Sum(1,2,3,4,5))
//}

//type Employee struct {
//ID   string
//Name string
//Age  int
//}
//func TestCreateEmployeeObj(t *testing.T) {
//e := Employee{"0", "Bob", 20}
//e1 := Employee{Name: "Mike", Age: 30}
//e2 := new(Employee)
//e2.ID = "2"
//e2.Name = "Rose"
//e2.Age = 22
//t.Log(e)
//t.Log(e1)
//t.Log(e1.ID)
//t.Log(e2)
//t.Logf("e is %T", e)
//t.Logf("e2 is %T", e2)
//}

// 实例的成员不会进行值复制
//func (e *Employee) String() string {
//fmt.Printf("Address is %x\n", unsafe.Pointer(&e.Name))
//return fmt.Sprintf("ID:%s-Name:%s-Age:%d", e.ID, e.Name, e.Age)
//}

// 实例的成员会进行值复制
//func (e Employee) String() string {
//fmt.Printf("Address is %x\n", unsafe.Pointer(&e.Name))
//return fmt.Sprintf("ID:%s-Name:%s-Age:%d", e.ID, e.Name, e.Age)
//}

//func TestStructOperations(t *testing.T) {
//e := &Employee{"0", "Bob", 20}
//fmt.Printf("Address is %x\n", unsafe.Pointer(&e.Name))
//t.Log(e.String())
//}

// interface ...

//type Programmer interface {
//WriteHelloWorld() string
//}

//type GoProgrammer struct {

//}

//func (g *GoProgrammer) WriteHelloWorld() string{
//return "fmt.Println(\"Hello World\")"
//}

//func TestClient(t *testing.T) {
//var p Programmer
//p = new(GoProgrammer)
//t.Log(p.WriteHelloWorld())
//}
