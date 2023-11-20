package goirc 

import (
    "fmt"
    "log"
    "bufio" 
    "net"
    "net/textproto"
    "strings"
    "syscall"
    "os"
    "os/signal"
)

const (
    PRIVMSG = iota
    NONE
)

type MSG struct {
    Username string
    chat_type int
    Message string
}

func (msg *MSG)IsPrivateMessage()(bool) {
    if msg.chat_type == PRIVMSG {
        return true
    }
    return false
}

type Client struct {
    address string
    port string
    password string
    server string
    username string
    nick string
    chat MSG 
    conn net.Conn
    reader *textproto.Reader
    disconnected bool
}

func handle_error(err error){
    if err != nil {
        log.Fatal(err)
    }
}

func send_data(conn net.Conn, msg string) {
    fmt.Fprintf(conn, "%s\n", msg)
}

func parse_message(message string) (MSG) {
    var chat MSG 
    if strings.Contains(message, "PRIVMSG") {
        chat.chat_type = PRIVMSG
    } else {
        chat.chat_type = NONE
    }


    index := strings.Index(message, " :")
    if index != -1 {
        chat.Message = message[index + 2:]
    } else {
        chat.Message = message
    }

    if chat.chat_type == PRIVMSG {
        index = strings.Index(message, "!")
        username := message[:index] 
        chat.Username = username[1:]
    }


    return chat 
}

func Init(address string, port string, password string, username string, nick string)(Client) {
    var client Client
    client.address = address
    client.port = port
    client.password = password 
    client.username = username 
    client.nick = nick 
    return client
}

func (client *Client)Connect() {
    conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", client.address, client.port))
    handle_error(err)
    client.conn = conn
    client.Auth()
}

func (client *Client)Disconnect() {
    c := make(chan os.Signal, 1)
    signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        for range c {
            client.disconnected = true
            client.conn.Close()
            os.Exit(0)
        }
    }()
}

func (client *Client)Auth() {
    if len(client.password) > 0 {
        send_data(client.conn, "PASS " + client.password)
    }

    if len(client.username) > 0 {
        send_data(client.conn, "USER " + client.username)
    }

    if len(client.nick) > 0 {
        send_data(client.conn, "NICK " + client.nick)
    }
}

func (client *Client)Join(server string){
    send_data(client.conn, "JOIN #" + server)
    client.server = server

    reader := bufio.NewReader(client.conn)
    tp := textproto.NewReader(reader)
    client.reader = tp
}

func (client *Client)HandlePong(data string) {
    if client.disconnected {
        return
    }
    if strings.HasPrefix(data, "PING") {
        send_data(client.conn, fmt.Sprintf("PONG %s\n", strings.TrimPrefix(data, "PING ")))
        fmt.Println(fmt.Sprintf("PONG %s\n", strings.TrimPrefix(data, "PING ")))
    }
}

func (client *Client)GetData()(MSG) {
    if !client.disconnected {
        data, err := client.reader.ReadLine()
        handle_error(err)
        client.HandlePong(data)
        fmt.Println(data)
        msg := parse_message(data)
        client.chat = msg
        return msg 
    }
    return MSG{} 
}

func (client *Client)Say(msg string) {
    send_data(client.conn, fmt.Sprintf("PRIVMSG #%s :%s\n", client.server, msg))
}
