package shadowaead

import (
	"fmt"
	"sync"
	"encoding/hex"
	"net"
	"testing"

    "bytes"
    "crypto/rand"
)

func CheckConn(a net.Conn, b net.Conn) bool {
    payload1 := [1024]byte{}
    payload2 := [1024]byte{}
    rand.Reader.Read(payload1[:])
    rand.Reader.Read(payload2[:])

    result1 := [1024]byte{}
    result2 := [1024]byte{}
    wg := sync.WaitGroup{}
    wg.Add(2)
    go func() {
        a.Write(payload1[:])
        a.Read(result2[:])
        wg.Done()
    }()
    go func() {
        b.Read(result1[:])
        b.Write(payload2[:])
        wg.Done()
    }()
    wg.Wait()

    //fmt.Println(payload1[:], result1[:])
    if !bytes.Equal(payload1[:], result1[:]) || !bytes.Equal(payload2[:], result2[:]) {
        return false
    }
    return true
}

func TestCipher(t *testing.T) {
	key, _ := hex.DecodeString("ab06c33b91cb5fe9a386d1de81f8cbcc")

	aead, err := AESGCM(key)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	//fmt.Println(aead)

    l, err := net.Listen("tcp", "127.0.0.1:1341")
    if err != nil {
        t.Fatal(err)
    }

	wg := sync.WaitGroup{}
    wg.Add(1)
    var conn1, conn2 net.Conn
    go func() {
        conn1, err = l.Accept()
        if err != nil {
        	fmt.Println(err)
        	panic(err)
        }
        wg.Done()
    }()
    conn2, err = net.Dial("tcp", "127.0.0.1:1341")
    if err != nil {
    	fmt.Println(err)
    	panic(err)
    }
    conn2 = NewConn(conn2, aead)

	conn2.Write([]byte("12345678\r\n"))

    wg.Wait()

    conn1 = NewConn(conn1, aead)

    buf := [12]byte{}
    n, err := conn1.Read(buf[:])
    if err!=nil {
    	fmt.Println(n, err)
    	t.Fail()
    }
    

    //if !CheckConn(conn1, conn2) {
    //    t.Fail()
    //}
}

