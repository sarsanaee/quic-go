package main

import (
	rand "crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"math"
	"math/big"
	mrand "math/rand"
	"time"
	"os"
        "strconv"


	quic "github.com/lucas-clemente/quic-go"
)

const addr = "10.254.254.1:4242"
const callpersec = 500
const message = "foobar"

// We start a server echoing data on the first stream the client opens,
// then connect with a client, send the message, and wait for its receipt.
func main() {
	//go func() { log.Fatal(echoServer()) }()

        args := os.Args[1:]
        rate_int, err := strconv.Atoi(args[0])
	var rate float64 = float64(rate_int)


	mrand.Seed(time.Now().UnixNano()) //generating a new seed but
	//I guess it should be a number in all the experiments.
	//fmt.Println(nextTime(100.0))

	err = clientMain(rate)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

func nextTime(rate float64) float64 {
	return -1 * math.Log(1.0-mrand.Float64()) / rate
}

func clientWrite(stream quic.Stream, rate float64, err error) error {
	var my_random_number float64 = nextTime(rate) * 1000000
	var my_random_int int = int(my_random_number)
	var int_message int64 = time.Now().UnixNano()
	byte_message := make([]byte, 8)
	for true {

		time.Sleep(time.Microsecond * time.Duration(my_random_int))
		int_message = time.Now().UnixNano()
		binary.LittleEndian.PutUint64(byte_message, uint64(int_message))
		// _, err = stream.Write([]byte(byte_message))
		//fmt.Println("Send", byte_message, int_message)
		_, err = stream.Write(byte_message)

		if err != nil {
			return err
		}

	}
	return err
}

func clientRead(stream quic.Stream, err error) error {
	buf := make([]byte, 8) //len(message))
	for true {
		_, err = io.ReadFull(stream, buf)
		now := time.Now().UnixNano()

		if err != nil {
			return err
		}
		last := int64(binary.LittleEndian.Uint64(buf))
		fmt.Println((now - last) / 1000)
	}
	return err
}

func clientMain(rate float64) error {
	session, err := quic.DialAddr(addr, &tls.Config{InsecureSkipVerify: true}, nil)
	if err != nil {
		return err
	}

	stream, err := session.OpenStreamSync()
	if err != nil {
		return err
	}
	go clientRead(stream, err)

	clientWrite(stream, rate , err)

	// for true {
	// 	fmt.Printf("Client: Sending '%s'\n", message)
	// 	_, err = stream.Write([]byte(message))
	// 	if err != nil {
	// 		return err
	// 	}

	// 	buf := make([]byte, len(message))
	// 	_, err = io.ReadFull(stream, buf)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	fmt.Printf("Client: Got '%s'\n", buf)
	// }

	return nil
}

// A wrapper for io.Writer that also logs the message.
type loggingWriter struct{ io.Writer }

func (w loggingWriter) Write(b []byte) (int, error) {
	//fmt.Printf("Server: Got '%s'\n", string(b))
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

