package AxmorChat


import (
	"net"
	"sync"
)

type OnlineClientList struct {
	sync.Mutex
	List map [net.Conn]*Client

}
//Добавление нового клиента в ассоциативный массив
func (ocl *OnlineClientList) Add(conn net.Conn,client *Client){
	ocl.Lock()
	defer ocl.Unlock()
	ocl.List[conn] = client
}

//Получение клиента по ключу
func (ocl *OnlineClientList) Get(conn net.Conn) *Client{
	ocl.Lock()
	defer ocl.Unlock()
	return  ocl.List[conn]
}

//Удаление клиента из массива по ключу
func (ocl *OnlineClientList) Delete(conn net.Conn) {
	ocl.Lock()
	defer ocl.Unlock()
	delete(ocl.List, conn)
}