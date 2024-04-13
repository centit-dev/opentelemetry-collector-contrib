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

	host, port, err = extractHostAndPort("mysql://smartobserv_mysql.smartobserv-local/dbname")
	assert.Nil(t, err)
	assert.Equal(t, "smartobserv_mysql.smartobserv-local", host)
	assert.Equal(t, "", port)
}
