package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("creating new client returns correct instance", func(t *testing.T) {
		var in io.ReadCloser
		var out io.Writer
		address := "otus.ru:80"
		timeout := time.Duration(42)
		var conn net.Conn

		expectedClient := &Client{
			address: address,
			timeout: timeout,
			in:      in,
			out:     out,
			conn:    conn,
		}
		client := NewTelnetClient(address, timeout, in, out)

		require.Equal(t, expectedClient, client)
	})

	t.Run("return error if connection failed", func(t *testing.T) {
		var in io.ReadCloser
		var out io.Writer
		address := "otus.ru"
		timeout := time.Duration(42)

		client := NewTelnetClient(address, timeout, in, out)
		err := client.Connect()

		fmt.Println(err)
		require.Equal(t, "connection error: dial tcp: address otus.ru: missing port in address", err.Error())
	})

	t.Run("return error if connection failed", func(t *testing.T) {
		var in io.ReadCloser
		var out io.Writer
		address := "otus.ru:80"
		timeout := time.Duration(42)

		client := NewTelnetClient(address, timeout, in, out)
		err := client.Connect()

		fmt.Println(err)
		require.Equal(t, "connection error: dial tcp: i/o timeout", err.Error())
	})

	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, ioutil.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})
}
