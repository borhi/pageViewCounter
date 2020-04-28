package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"testing"
)

func Test_Main(t *testing.T) {
	port := chooseRandomUnusedPort()
	os.Args = []string{"test", "-port=" + strconv.Itoa(port)}
	callsCount := 100
	go func() {
		main()
	}()

	ch := make(chan http.Response)
	for i := 0; i < callsCount; i++ {
		go func() {
			r, err := http.Get(fmt.Sprintf("http://localhost:%d", port))
			require.NoError(t, err)
			ch <- *r
		}()
	}

	for i := 0; i < callsCount; i++ {
		<-ch
	}

	//get last result
	response, err := http.Get(fmt.Sprintf("http://localhost:%d", port))
	require.NoError(t, err)
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		require.NoError(t, err)
	}
	err = response.Body.Close()
	require.NoError(t, err)
	result, err := strconv.Atoi(string(bodyBytes))
	require.NoError(t, err)

	assert.Equal(t, callsCount+1, result)
}

func chooseRandomUnusedPort() (port int) {
	for i := 0; i < 10; i++ {
		port = 40000 + int(rand.Int31n(10000))
		if ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port)); err == nil {
			_ = ln.Close()
			break
		}
	}
	return port
}
