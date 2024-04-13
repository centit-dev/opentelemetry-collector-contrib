package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractHostAndPort(t *testing.T) {
	host, port, err := extractHostAndPort("mysql://user:pass@localhost:3306/dbname")
	assert.Nil(t, err)
	assert.Equal(t, "localhost", host)
	assert.Equal(t, "3306", port)
}
