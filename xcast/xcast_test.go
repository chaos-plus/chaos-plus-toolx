package xcast

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSnakeToCamel(t *testing.T) {
	require.Equal(t, "HelloWorld", SnakeToCamel("hello_world"))
	require.Equal(t, "HelloWorld", SnakeToCamel("Hello_World"))
	require.Equal(t, "HelloWorld", SnakeToCamel("helloWorld"))
	require.Equal(t, "Hello", SnakeToCamel("hello"))
	require.Equal(t, "A", SnakeToCamel("a"))
	require.Equal(t, "", SnakeToCamel("_"))
	require.Equal(t, "HelloWorld", SnakeToCamel("hello_world"))
	require.Equal(t, "HelloWorldABC", SnakeToCamel("hello_World_A_B_C"))
	require.Equal(t, "HelloWorldABC", SnakeToCamel("hello_World_ABC"))
}

func TestCamelToSnake(t *testing.T) {
	require.Equal(t, "hello_world", CamelToSnake("HelloWorld"))
	require.Equal(t, "hello_world_abc", CamelToSnake("HelloWorldAbc"))
	require.Equal(t, "hello", CamelToSnake("Hello"))
	require.Equal(t, "", CamelToSnake(""))
	require.Equal(t, "a", CamelToSnake("A"))
	require.Equal(t, "hello_world_abc", CamelToSnake("HelloWorldABC"))
	require.Equal(t, "hello_world_abc", CamelToSnake("HelloWorldAbc"))
}

func TestCast(t *testing.T) {

	require.Equal(t, false, ToBool("0"))
	require.Equal(t, true, ToBool("1"))

	require.Equal(t, float64(1.0), ToFloat64("1"))
	require.Equal(t, float32(1.0), ToFloat32("1"))

	require.Equal(t, 1, ToInt("1"))
	require.Equal(t, int8(1), ToInt8("1"))
	require.Equal(t, int16(1), ToInt16("1"))
	require.Equal(t, int32(1), ToInt32("1"))
	require.Equal(t, int64(1), ToInt64("1"))

	require.Equal(t, uint(1), ToUint("1"))
	require.Equal(t, uint8(1), ToUint8("1"))
	require.Equal(t, uint16(1), ToUint16("1"))
	require.Equal(t, uint32(1), ToUint32("1"))
	require.Equal(t, uint64(1), ToUint64("1"))

	require.Equal(t, "1", ToString("1"))
	require.Equal(t, "1", ToString(1))

}
