package object

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_mapKey(t *testing.T) {
	s1, s2 := &Stringer{Value: "s"}, &Stringer{Value: "s"}
	assert.Equal(t, s1.MapKey(), s2.MapKey())
	b1, b2 := &Boolean{Value: true}, &Boolean{Value: true}
	assert.Equal(t, b1.MapKey(), b2.MapKey())
	i1, i2 := &Integer{Value: 1}, &Integer{Value: 1}
	assert.Equal(t, i1.MapKey(), i2.MapKey())
}
