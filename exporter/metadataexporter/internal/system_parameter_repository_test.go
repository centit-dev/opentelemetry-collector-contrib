package internal

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func DisabledTestFindByCode(t *testing.T) {
	client := createTestClient()
	repository := CreateSystemParameterRepository(client)
	value, _ := repository.FindByCode(context.Background(), metadataKeyBlacklist)
	assert.NotEmpty(t, value)
}
