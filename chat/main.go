package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/juancwu/Gorangu/chat/views"
	"github.com/labstack/echo/v4"
	"github.com/matoous/go-nanoid/v2"
)

type Message struct {
	Content string `json:"chat_message"`
}

type Client struct {
	conn *websocket.Conn
    clientId string
	data chan string
}

type Broadcast struct {
    data string
    clientId string
}

type Room struct {
    clients map[string]*Client
    broadcast chan Broadcast
    register chan *Client
    unregister chan string
}

var (
	upgrader = websocket.Upgrader{}
    rooms = make(map[string]Room)
)

func (r Room) run(ctx echo.Context, roomId string) {
    ctx.Logger().Info(fmt.Sprintf("Running room: %s\n", roomId))
    for {
        select {
        case b := <- r.broadcast:
            for _, c := range r.clients {
                c.data <- b.data
            }
        case client := <- r.register:
            r.clients[client.clientId] = client
        case clientId := <- r.unregister:
            _, ok := r.clients[clientId]
            if ok {
                delete(r.clients, clientId)
            }

            fmt.Printf("Number of clients in room (%s): %d\n", roomId, len(r.clients))
            
            if len(r.clients) == 0 {
                delete(rooms, roomId)
            }
        }
    }
}

func (c Client) readMessage(ctx echo.Context, roomId string) {
	err := c.conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		ctx.Logger().Error(err)
		return
	}

	c.conn.SetPongHandler(func(appData string) error {
		err := c.conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		if err != nil {
			ctx.Logger().Error(err)
			return err
		}
		fmt.Println("pong")
		return nil
	})

    room, ok := rooms[roomId]
    if !ok {
		fmt.Println("Closing websocket.")
		c.conn.Close()
        return
    }

	defer func() {
		fmt.Println("Closing websocket...")
        room.unregister <- c.clientId
		c.conn.Close()
	}()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			ctx.Logger().Error(err)
            return
		}
		var message Message
		err = json.Unmarshal(msg, &message)
		if err != nil {
			ctx.Logger().Error(err)
			continue
		}
		fmt.Printf("From client: %s\n", message.Content)
        b := Broadcast{message.Content, c.clientId}
        room.broadcast <- b
	}
}

func (c Client) writeMessage(ctx echo.Context) {
	defer c.conn.Close()

	ticker := time.NewTicker(time.Second * 9)
	for {
		select {
		case text, ok := <-c.data:
			if !ok {
				return
			}
            component := views.Message(text)
            buffer := &bytes.Buffer{}
            component.Render(context.Background(), buffer)
			fmt.Printf("Message to send: %s\n", text)
            err := c.conn.WriteMessage(websocket.TextMessage, buffer.Bytes())
			if err != nil {
				ctx.Logger().Error(err)
				return
			}
		case <-ticker.C:
			err := c.conn.WriteMessage(websocket.PingMessage, []byte(""))
			if err != nil {
				ctx.Logger().Error(err)
				return
			}
			fmt.Println("Ping")
		}
	}
}

func main() {
	fmt.Println("chat")

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
        roomId := gonanoid.Must(6)
        c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/chat/%s", roomId))
		return nil
	})
    e.GET("/chat/:roomId", func (c echo.Context) error {
        roomId := c.Param("roomId")
		views.Index(roomId).Render(context.Background(), c.Response().Writer)
		return nil
    })
    e.GET("/chatroom/:roomId", func(c echo.Context) error {
        roomId := c.Param("roomId")
        room, ok := rooms[roomId]
        if !ok {
            room = Room{
                make(map[string]*Client),
                make(chan Broadcast),
                make(chan *Client),
                make(chan string),
            }
            rooms[roomId] = room
        }

		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}

        clientId := gonanoid.Must(12)
		client := Client{ws, clientId, make(chan string)}

		go client.readMessage(c, roomId)
        go client.writeMessage(c)
        go room.run(c, roomId)

        room.register <- &client

		return nil
	})
    e.Static("/static", "static")

	e.Logger.Fatal(e.Start(":5173"))
}
