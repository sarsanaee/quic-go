package main

import (
	rand "crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"math"
	"math/big"
	mrand "math/rand"
	"time"

	quic "github.com/lucas-clemente/quic-go"
)

const addr = "0.0.0.0:4242"

const message = "foobar"

// We start a server echoing data on the first stream the client opens,
// then connect with a client, send the message, and wait for its receipt.
func main() {
	//go func() { log.Fatal(echoServer()) }()

	//fmt.Println(time.Microsecond * nextTime(100.0) * 1000000)

	//time.Sleep(nextTime(100.0) * 1000000)

	//fmt.Println(nextTime(1000.0))
	echoServer()
	//err := clientMain()
	// err := echoServer()
	// if err != nil {
	// 	log.Fatal(err)
	// 	panic(err)
	// }
}

func nextTime(rate float64) float64 {
	return -1 * math.Log(1.0-mrand.Float64()) / rate
}

func echoStream(sess quic.Session) {
	for true {
		stream, err := sess.AcceptStream()
		if err != nil {
			panic(err)
			break
		}
		// Echo through the loggingWriter
		fmt.Println("before IO")
		_, err = io.Copy(loggingWriter{stream}, stream)
		fmt.Println("after IO")
		//return err
	}
}

// Start a server that echos all data on the first stream opened by the client
func echoServer() {

	listener, err := quic.ListenAddr(addr, generateTLSConfig(), nil)
	fmt.Println("after listener")

	for true {

		if err != nil {
			//return err
			break
		}

		fmt.Println("before accept")
		sess, err := listener.Accept()
		fmt.Println("after accept")
		if err != nil {
			//return err
			break
		}

		go echoStream(sess)

	}
	return
	//return err
}

// A wrapper for io.Writer that also logs the message.
type loggingWriter struct{ io.Writer }

func (w loggingWriter) Write(b []byte) (int, error) {
	//fmt.Printf("Server: Got '%s'\n", string(b))##########################
	time.Sleep(time.Microsecond * 10) //time.Duration(nextTime(.0)))
	return w.Writer.Write(b)
}

// Setup a bare-bones TLS config for the server
func generateTLSConfig() *tls.Config {
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
