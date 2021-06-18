package main

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTelent(t *testing.T) {
	cmd := exec.Command("go", "run", "telnet.go", "www.qq.com", "80")
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	err := cmd.Start()
	assert.Nil(t, err, err)
	stdin.Write([]byte("GET / HTTP/1.1\r\nHost: www.qq.com\r\n\r\n"))
	buff := make([]byte, 4096)
	bytes, err := stdout.Read(buff)
	assert.Nil(t, err, err)
	assert.Greater(t, bytes, 0, "bytes should > 0")
	t.Logf("%s", string(buff[:bytes]))
	assert.Equal(t, string(buff[:8]), "HTTP/1.1")
}
