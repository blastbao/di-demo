package main


import (
	"database/sql"
	"fmt"
	"os"
	_ "github.com/go-sql-driver/mysql"
	"github.com/blastbao/di-demo/demo"
	"github.com/blastbao/di-demo/di"
)


func main() {

	db, err := sql.Open("mysql", "root:root@tcp(localhost)/sampledb")
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		os.Exit(1)
	}


	container := di.NewContainer()
	container.SetSingleton("db", db)
	container.SetPrototype("b", func() (interface{}, error) {
		return demo.NewB(), nil
	})

	a := demo.NewA()
	if err := container.Ensure(a); err != nil {
		fmt.Println(err)
		return
	}
	// 打印指针，确保单例和实例的指针地址
	fmt.Printf("db0: %p\ndb1: %p\nb0: %p\nb1: %p\n", a.Db0, a.Db1, &a.B0, &a.B1)
	fmt.Printf("db0: %v\ndb1: %v\nb0: %v\nb1: %v\n", a.Db0, a.Db1, &a.B0, &a.B1)
}
