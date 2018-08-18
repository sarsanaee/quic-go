package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"math/big"
	mrand "math/rand"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	quic "github.com/lucas-clemente/quic-go"
)

const addr = "localhost:4242"

const message = "foobar"

const filename = "/home/alireza/quic_results/quic/"

var total_rcv int64
var syncFlag bool

// We start a server echoing data on the first stream the client opens,
// then connect with a client, send the message, and wait for its receipt.
func main() {

	tPtr := flag.String("type", "client", "application type")
	cmdRateIntPtr := flag.Float64("rate", 400000, "change rate of message reading")
	cmdPortPtr := flag.String("port", ":9090", "port to listen")
	clientSizePtr := flag.Int("size", 20, "number of clients")
	serverIPPtr := flag.String("ip", "10.254.254.239", "server_ip")
	expTimePtr := flag.Int("time", 5, "Experiment time")

	flag.Parse()

	t := *tPtr
	cmdRateInt := *cmdRateIntPtr
	cmdPort := *cmdPortPtr
	clientSize := *clientSizePtr
	serverIP := *serverIPPtr
	expTime := *expTimePtr

	if t == "server" {
		server(cmdPort, serverIP)
	} else if t == "client" {

		syncFlag = true
		t1 := time.Now()

		for i := 0; i < clientSize; i++ {
			go clientMain(cmdRateInt, serverIP, cmdPort)
		}
		// <-make(chan bool) // infinite wait.
		<-time.After(time.Second * time.Duration(expTime))
		syncFlag = false
		<-time.After(time.Second * time.Duration(expTime))
		// fmt.Println("total exchanged:", total_rcv, "\nthroughput:",
		// 	total_rcv*1000000000/time.Now().Sub(t1).Nanoseconds(), "call/sec")
		writeThroughput(total_rcv*1000000000/time.Now().Sub(t1).Nanoseconds(), cmdRateInt)
	}
	// go func() { log.Fatal(echoServer()) }()

	// err := clientMain()
	// if err != nil {
	// 	panic(err)
	// }
}

func echo(sess quic.Session) {

	stream, err := sess.AcceptStream()
	buf := make([]byte, 8)

	if err != nil {
		panic(err)
	}

	// Echo through the loggingWriter
	//_, err = io.Copy(loggingWriter{stream}, stream)

	for {
		_, err := io.ReadFull(stream, buf)
		if err != nil {
			return
		}
		//fmt.Println("recv")
		time.Sleep(time.Microsecond * 10)
		stream.Write(buf)
	}

	// for {
	// 	_, err := io.ReadFull(conn, buf)
	// 	if err != nil {
	// 		return
	// 	}
	// 	//fmt.Println("recv")
	// 	time.Sleep(time.Microsecond * 10)
	// 	conn.Write(buf)
	// }
}

// Start a server that echos all data on the first stream opened by the client
func server(cmdPort string, serverIP string) error {

	//fmt.Println(serverIP + cmdPort)

	listener, err := quic.ListenAddr(serverIP+cmdPort, generateTLSConfig(), nil)
	if err != nil {

		return err
	}

	for {
		sess, err := listener.Accept()
		if err != nil {
			return err
		}
		go echo(sess)

	}

	return err
}

func clientMain(cmdRateInt float64, serverIP string, cmdPort string) error {
	session, err := quic.DialAddr(serverIP+cmdPort, &tls.Config{InsecureSkipVerify: true}, nil)
	if err != nil {
		return err
	}

	stream, err := session.OpenStreamSync()
	if err != nil {
		return err
	}

	byte_message := make([]byte, 8)

	// fmt.Printf("Client: Sending '%s'\n", message)
	// _, err = stream.Write([]byte(message))
	// if err != nil {
	// 	return err
	// }
	defer stream.Close()

	go func(stream quic.Stream) {

		var rtt int64
		var latency string

		buf := make([]byte, 8)
		for {
			if !syncFlag {
				myPrint(latency)
				break
			}

			_, err := io.ReadFull(stream, buf)
			if err != nil {
				break
			}
			int_message := int64(binary.LittleEndian.Uint64(buf))
			t2 := time.Unix(0, int_message)
			rtt = (time.Now().UnixNano() - t2.UnixNano()) / 1000
			latency = latency + strconv.FormatInt(rtt, 10) + "\n"
			//fmt.Print(latency)

			//fmt.Println((time.Now().UnixNano() - t2.UnixNano()) / 1000)
			atomic.AddInt64(&total_rcv, 1)
		}
		return
	}(stream)

	for {

		wait := time.Microsecond * time.Duration(nextTime(cmdRateInt)*1000000)
		if wait > 0 {
			time.Sleep(wait)
			//fmt.Println("WAIT", wait)
		}
		int_message := time.Now().UnixNano()
		binary.LittleEndian.PutUint64(byte_message, uint64(int_message))
		// _, err := conn.Write(byte_message)
		_, err = stream.Write(byte_message)

		if err != nil {
			log.Println("ERROR", err)
			return err
		}
	}

	// buf := make([]byte, len(message))
	// _, err = io.ReadFull(stream, buf)
	// if err != nil {
	// 	return err
	// }
	// fmt.Printf("Client: Got '%s'\n", buf)

	return nil
}

// A wrapper for io.Writer that also logs the message.
type loggingWriter struct{ io.Writer }

func (w loggingWriter) Write(b []byte) (int, error) {
	// fmt.Printf("Server: Got '%s'\n", string(b))
	time.Sleep(time.Microsecond * 10)
	return w.Writer.Write(b)
}

// Setup a bare-bones TLS config for the server
func generateTLSConfig() *tls.Config {
	//fmt.Println("SHITHIHIHTI")
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{Certificates: []tls.Certificate{tlsCert}}
}

func nextTime(rate float64) float64 {
	return -1 * math.Log(1.0-mrand.Float64()) / rate
}

func myPrint(latency_series string) {
	fmt.Print(latency_series)
}

func writeThroughput(throughput int64, rate float64) {
	curFileName := filename + strconv.Itoa(int(rate)) + "/throughputs.txt"
	f, err := os.OpenFile(curFileName, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		f, err := os.Create(curFileName)
		if err != nil {
			panic(err)
		}
		f.Close()
	}
	f, err = os.OpenFile(curFileName, os.O_APPEND|os.O_WRONLY, 0600)

	defer f.Close()

	_, err = f.WriteString(strconv.Itoa(int(throughput)) + "\n")

	if err != nil {
		panic(err)
	}
}
