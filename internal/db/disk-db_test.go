package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiskDb(t *testing.T) {
	db := diskDb("/tmp/kitproj/sim/db")

	assert.False(t, db.Put("foo", "bar"))
	assert.True(t, db.Put("foo", "bar"))

	assert.ElementsMatch(t, []any{"bar"}, db.List("foo"))
	assert.Equal(t, "bar", db.Get("foo"))

	assert.Nil(t, db.Get("bar"))

	assert.True(t, db.Delete("foo"))
	assert.False(t, db.Delete("foo"))

}
