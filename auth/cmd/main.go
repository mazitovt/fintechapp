package main

import (
	"github.com/mazitovt/fintechapp/auth/app"
	"log"
)

func main() {
	if err := app.Run("config"); err != nil {
		log.Fatal(err)
	}
}

//func Run1() error {
//
//	s := "mysecret"
//	sub := "tim@ya.ru"
//	m := NewManager(s)
//
//	tkn, err := m.NewAccessToken(sub, 10*time.Second)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Println(tkn)
//
//	actSub, err := m.Parse(tkn)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Println(actSub)
//	fmt.Println(sub == actSub)
//
//	p := []byte(tkn)
//	p[4] = 84
//	tkn = string(p)
//
//	actSub2, err := m.Parse(tkn)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Println(actSub2)
//
//	return nil
//}
