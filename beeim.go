package main

import (
	"log"
	"github.com/jaydenhe/BeeIM/s2c"
)

func main(){
	s2c_server := s2c.CreateServer()
	log.Printf("Running on \n")
	s2c_server.Start(":1114")
}
