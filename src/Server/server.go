package main

import (
	. "./Chat"
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

const byeMsg = "Всем пока! \n"
const maxLenHistory  = 10

var (
	onlineClientList = OnlineClientList {List: map[net.Conn]*Client{}}
	offlineClientList =  make ([]*Client,0)
	newConnection   = make(chan net.Conn)
	deadConnection  = make(chan net.Conn)
	sendList = make(chan net.Conn)
	sendAll = make(chan *Message)
	msgHistory = make ([]*Message,0)
)


func main() {

	ln, err := net.Listen("tcp", ":8080")
	logFatal(err)

	log.Println("Запущен сервер ", ln.Addr())

	defer ln.Close()
	/*
	Поток для подключения клиентов
	*/
	go func() {
		for {
			conn, err := ln.Accept()
			logFatal(err)

			reader := bufio.NewReader(conn)
			name, err := reader.ReadString('\n')
			name = strings.Trim(name,"\r\n")

			onlineClientList.Add(conn, NewClient(name, true))

			log.Printf("Клиент %s присоеденился", name)
			newConnection <- conn
		}
	}()
	/*
	Основной цикл чтения и записи от клиента
	*/
	for {
		select {
		case conn := <-newConnection:
			// Отправка истории сообщений
			log.Println("Отправка истории клиенту", conn)
			if len(msgHistory) != 0 {
				for _, msg := range msgHistory {
					sendMessage(conn, ReturnMessage(msg))
					time.Sleep(1 * time.Microsecond)
				}
			}
			//Запуск потока прослушивания клиента
			go ListenClient(conn)
		//Отключение клиента
		case conn := <-deadConnection:
			log.Printf("%s ушел в оффлайн", onlineClientList.Get(conn).Name)
			onlineClientList.Get(conn).Connected = false
			offlineClientList = append(offlineClientList, onlineClientList.Get(conn))
			onlineClientList.Delete(conn)

		//Отправка сообщения всем, кроме того клиента который его отправил
		case msg:= <- sendAll:

			onlineClientList.Lock()
			for item,client := range onlineClientList.List {
				if client != msg.Client {
					sendMessage(item, ReturnMessage(msg))
				}
			}
			onlineClientList.Unlock()

			log.Println("Запись сообщения в историю", ReturnMessage(msg))
			msgHistory = append(msgHistory, msg)
			//Проверка количества сообщений в истории и удаление старых
			if len(msgHistory) > maxLenHistory {
				msgHistory = msgHistory[1:maxLenHistory+1]
			}
			log.Println("Количество соощений в истории", len(msgHistory))

		//Запрос статуса клиентов чата
		case conn:= <-sendList:
			log.Println("Запрос статуса клиента")
			onlineClientList.Lock()
			for _,client := range onlineClientList.List{
				sendMessage(conn, StatusClient(client))
				time.Sleep(1*time.Microsecond)
			}
			onlineClientList.Unlock()
			for _,client := range offlineClientList{
				sendMessage(conn, StatusClient(client))
				time.Sleep(1*time.Microsecond)
			}
		}
	}
}

/*
Функция прослушивание соединения
*/
func ListenClient(conn net.Conn) {
	for {
		reader := bufio.NewReader(conn)
		message, err := reader.ReadString('\n')

		if err != nil {
			log.Println("Ошибка чтения данных:",err)
			break
		}

		if strings.Compare(strings.Trim(message,"\r\n"), "list") == 0 {
			sendList <- conn

		}else{

			if strings.Compare(strings.Trim(message, "\r\n"), "bye") == 0 {
				message = byeMsg
			}
			client := onlineClientList.Get(conn)
			log.Println("Новое сообщение от клиента ", client.Name)
			msg := Message{Client: client, Msg: message}
			log.Println("Новое сообщение ", ReturnMessage(&msg))
			sendAll <- &msg
		}
	}

	deadConnection <- conn
}

//Функция отправки сообщения
func sendMessage (conn net.Conn, msg string){
	_, err := fmt.Fprintf(conn, msg)

	if err != nil{
		log.Println("Ошибка отправки сообщения", err)
	}
}


