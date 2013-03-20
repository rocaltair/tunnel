package main

import (
	"fmt"
	"flag"
	"net"
	"crypto/rc4"
	"log"
	"strings"
	"strconv"
	"io"
)

const CHAN_SIZE = 20 
const BUFF_SIZE = 16

var From = flag.String("f", "127.0.0.1:6060", "endpoint where stream from,fmt:'host:port|:port|port'")
var To = flag.String("t", "127.0.0.1:5432", "endpoint where stream to,fmt:'host:port|:port|port'")
var Passwd = flag.String("p", "passwd", "your passwd")
var Debug = flag.Bool("d", false, "using debug mode")
var Server = flag.Bool("s", false, "as a server")

func NewEncoderDecoder(Passwd string, Server bool) (encoderFunc, decoderFunc func ([]byte) []byte, err error) {
	encoder, erre := rc4.NewCipher([]byte(Passwd))
	if erre != nil {
		return nil, nil, erre
	}
	decoder, errd := rc4.NewCipher([]byte(Passwd))
	if errd != nil {
		return nil, nil, errd
	}
	
	encoderFunc = func(in []byte) []byte {
		buff := make([]byte, len(in))
		encoder.XORKeyStream(buff, in)
		return buff
	}
	decoderFunc = func(in []byte) []byte {
		buff := make([]byte, len(in))
		decoder.XORKeyStream(buff, in)
		return buff
	}
	if Server {
		return decoderFunc, encoderFunc, nil
	}
	return encoderFunc, decoderFunc, nil
}

func MakeChan(conn net.Conn) <- chan string {
	out := make(chan string, CHAN_SIZE)
	go func(conn net.Conn, out chan string) {
		defer func () {
			conn.Close()
		}()
		buff := make([]byte, BUFF_SIZE)
		var connReader io.Reader = conn
		for {
			cnt, err := connReader.Read(buff)
			if err != nil {
				if *Debug {
					log.Println(err)
				}
				close(out)
				return 
			}
			out <- string(buff[0:cnt])
		}
	}(conn, out)
	return out
}

func HandleConn(connFrom net.Conn, ToAddr, Password string) {
	defer connFrom.Close()
	connTo, errTo := net.Dial("tcp", ToAddr)
	if errTo != nil {
		log.Println(errTo)	
		return 
	}
	defer connTo.Close()

	fchan := MakeChan(connFrom)
	tchan := MakeChan(connTo)

	encoder, decoder, errc := NewEncoderDecoder(Password, *Server)
	if errc != nil {
		log.Println("create encoder decoder error:", errc)
		return 
	}

	var connFromWriter io.Writer = connFrom
	var connToWriter io.Writer = connTo
	for {
		select {
		case data, e := <- fchan:
			if !e {
				if *Debug {
					log.Println("fchan closed", e)
				}
				return 
			}
			connToWriter.Write(encoder([]byte(data)))
		case data, e := <- tchan:
			if !e {
				if *Debug {
					log.Println("tchan closed", e)
				}
				return 
			}
			connFromWriter.Write(decoder([]byte(data)))
		}
	}
}

func StartTunnel(FromAddr, ToAddr, Password string) {
	l, err := net.Listen("tcp", FromAddr)
	if err != nil {
		log.Fatal("listen err:", err)
	}
	for {
		conn, errl := l.Accept()
		if errl != nil {
			if *Debug {
				log.Println("accept error", errl)
			}
			continue
		}
		go HandleConn(conn, ToAddr, Password)
	}
}

func main() {
	flag.Parse()
	FromAddr := strings.TrimSpace(*From)
	if _, err := strconv.Atoi(FromAddr); err == nil {
		FromAddr = ":" + FromAddr
	}
	ToAddr := strings.TrimSpace(*To)
	if _, err := strconv.Atoi(ToAddr); err == nil {
		ToAddr = ":" + ToAddr
	}
	fmt.Printf("connect from '%s' to '%s', using passwd:'%s'\n", FromAddr, ToAddr, *Passwd)
	StartTunnel(FromAddr, ToAddr, *Passwd)
}

