package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInteralValidate(t *testing.T) {
	t.Parallel()

	p := Provider().(*schema.Provider)

	assert.Nil(t, p.InternalValidate())
}
