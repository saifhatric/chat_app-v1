package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/gorilla/websocket"
)


var wsChan=make(chan Payload)
var clients=make(map[*websocket.Conn]string)
type websocketConnection struct{
	*websocket.Conn
}

type Response struct{
	Action			string		`json:"action"`
	Username 		string		`json:"username"`
	Msg				string		`json:"message"`
	ConnectedUsers []string		`json:"connected_users"`
}
type Payload struct{
	Action		string		`json:"action"`
	Username 	string		`json:"username"`
	Msg			string		`json:"message"`
	Conn		websocketConnection	

}
var upgradeConnection=websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {return true},
}


func WsEndpoint(w http.ResponseWriter, r *http.Request){
	ws,err:=upgradeConnection.Upgrade(w,r,nil)
	if err!=nil{
		log.Fatal(err.Error())
	}
	var res Response
	res.Msg="Welcome !"
	conn:=websocketConnection{ws}
	fmt.Println("new websocket connection started")
	if err:=ws.WriteJSON(res);err!=nil{
		log.Fatalf("could not send response err: %v\n",err)
	}
	go socketListener(&conn)
}


func socketListener(conn *websocketConnection){
	defer func(){
		if r:=recover();r!=nil{
			log.Printf("Error : %v\n",r)
		}
		var payload Payload
	for{
		err:=conn.ReadJSON(&payload)
		if err!=nil{
			log.Fatal(err)
		}else{
			payload.Conn= *conn
			wsChan <-payload

		}
	}
}()
}

func ChanListener(){
	var res Response
	for {
		e:= <-wsChan
		switch e.Action{
			case "username":
				clients[e.Conn.Conn]=e.Username
				res.Action="list_users"
				users:=GetUsers()
				res.ConnectedUsers=users
				broadcastToAll(res)
				fmt.Println(res)
			case "left":
				res.Action="list_users"
				delete(clients,e.Conn.Conn)
				users:=GetUsers()
				res.ConnectedUsers=users
				broadcastToAll(res)
			case "broadcast":
				res.Action="broadcast"
				res.Msg=fmt.Sprintf("%s : %s",e.Username,e.Msg)
				broadcastToAll(res)
	}
	}
}

func GetUsers()[]string{
	var users []string
	for _,x:=range clients{
		if x!=""{
			users = append(users, x)
		}
	}
	sort.Strings(users)
	return users
}

func broadcastToAll(res Response){
	for x:= range clients{
		err:=x.WriteJSON(res)
		if err!=nil{
			log.Fatal(err)
			_=x.Close()
			delete(clients,x)
		}
	}
}


