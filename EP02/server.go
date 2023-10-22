package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type client chan<- string

type User struct {
	ip       string
	nickname string
	channel  chan<- string
	conn     net.Conn
}

type Message struct {
	sender  string
	destiny string
	text    string
}

type Command struct {
	command string
	args    string
}

var (
	entering = make(chan User)
	leaving  = make(chan User)
	editing  = make(chan User)
	messages = make(chan Message)
	commands = make(chan Command)
)

func broadcaster() {
	clients := make(map[User]bool) // todos os clientes conectados
	users := make(map[string]User)
	for {
		select {
		case msg := <-messages:
			// broadcast de mensagens. Envio para todos
			if msg.destiny == "all" {
				for i := range users {
					user := users[i]
					if strings.ToUpper(user.nickname) != strings.ToUpper(msg.sender) {
						user.channel <- msg.text
					}
				}
			} else {
				for i := range users {
					user := users[i]

					if strings.ToUpper(user.nickname) == strings.ToUpper(msg.destiny) {
						user.channel <- msg.text
					}
				}
			}

		case cli := <-entering:
			clients[cli] = true
			users[cli.ip] = cli

		case edited_cli := <-editing:
			users[edited_cli.ip] = edited_cli
			users[edited_cli.ip].channel <- "Seu username agora é " + edited_cli.nickname

		case user := <-leaving:
			delete(clients, user)
			delete(users, user.ip)
			close(user.channel)
			user.conn.Close()
		}

	}
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string)
	user := User{ip: conn.RemoteAddr().String(), nickname: conn.RemoteAddr().String(), channel: ch, conn: conn}

	go clientWriter(conn, ch)

	user.channel <- "Seu username é " + user.nickname
	messages <- Message{sender: user.nickname, text: user.nickname + " chegou!", destiny: "all"}
	entering <- user

	input := bufio.NewScanner(conn)
	for input.Scan() {
		if strings.HasPrefix(input.Text(), "/changenickname") {
			new_nick := strings.Split(input.Text(), "/changenickname ")[1]
			messages <- Message{sender: "server", text: user.nickname + " mudou de nome para " + new_nick, destiny: "all"}
			user.nickname = new_nick
			editing <- User{ip: user.ip, nickname: new_nick, channel: user.channel}

		} else if strings.HasPrefix(input.Text(), "/msg") {
			message_split := strings.Split(input.Text(), " ")

			destiny := message_split[1]
			message := strings.Join(message_split[2:], " ")

			messages <- Message{sender: user.nickname, text: "(private) " + user.nickname + ": " + message, destiny: destiny}
		} else if strings.HasPrefix(input.Text(), "/quit") || strings.HasPrefix(input.Text(), "/q") {
			messages <- Message{sender: "server", text: user.nickname + " deixou a sala.", destiny: "all"}
			leaving <- user
		} else if strings.HasPrefix(input.Text(), "/checkip") || strings.HasPrefix(input.Text(), "/ip") {
			messages <- Message{sender: "server", text: user.ip, destiny: user.ip}
		} else {
			messages <- Message{sender: user.nickname, text: user.nickname + ": " + input.Text(), destiny: "all"}
		}
	}

	conn.Close()
}

func main() {
	fmt.Println("Iniciando servidor...")
	listener, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		log.Fatal(err)
	}
	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}
