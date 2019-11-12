package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)


func main() {
	fmt.Println ("Enter your name")
	reader := bufio.NewReader(os.Stdin)
	name, _ := reader.ReadString('\n')
	fmt.Println ("Введите адрес сервера, например 127.0.0.1:8080")
	adr,_ := reader.ReadString('\n')
	adr = strings.Trim(adr, "\r\n")

	// Подключаемся к сокету
	conn, err := net.Dial("tcp", adr)
	if err != nil{
		fmt.Println("Сервер не доступен, попробуйте позже ")
	}
	defer os.Exit(1)

	fmt.Println("Присоеденились к серверу", conn.LocalAddr())

	// Отправляем в socket
	_, err = fmt.Fprintf(conn, name+"\r\n")
	if err != nil{
		fmt.Println(err)
	}

	go read(conn)
	write(conn)
}

func read(conn net.Conn) {
	for {
		reader := bufio.NewReader(conn)
		message, err := reader.ReadString('\n')

		if err !=nil {
			conn.Close()
			fmt.Println("Сервер не доступен, попробуйте позже.")
			os.Exit(1)
		}
		fmt.Print(message)
	}
}

func write (conn net.Conn){
	for{
		// Чтение входных данных от stdin
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')

		// Отправляем в socket
		_, err := fmt.Fprintf(conn, message+"\r\n")
		if err != nil{
			fmt.Println(err)
		}
		if strings.Compare(strings.Trim(message, "\r\n"), "bye") == 0 {
			break
		}
	}
}



