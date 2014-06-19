package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var target = Target{
	host:     "localhost",
	user:     "user",
	password: "password",
	path:     "/var/www/app",
}

func Test_getAddress(t *testing.T) {
	assert.Equal(t, target.getAddress(), "localhost:22")
}
