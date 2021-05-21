package app

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContextBehavior(t *testing.T) {
	ctx1 := context.Background()

	ctx2, cancel2 := context.WithCancel(ctx1)
	defer cancel2()

	ctx3, cancel3 := context.WithCancel(ctx2)
	defer cancel3()

	ctx4, cancel4 := context.WithCancel(ctx2)
	defer cancel4()

	ctx5, cancel5 := context.WithCancel(ctx4)
	defer cancel5()

	cancel4()

	require.NoError(t, ctx1.Err())
	require.NoError(t, ctx2.Err())
	require.NoError(t, ctx3.Err())
	require.Error(t, ctx4.Err())
	require.Error(t, ctx5.Err())
}

func TestContextDeviceID(t *testing.T) {
	div := Div()
	disp := NewClientTester(div)
	defer disp.Close()

	ctx := makeContext(div)
	id := ctx.DeviceID()
	require.NotZero(t, id)

	id2 := ctx.DeviceID()
	require.Equal(t, id, id2)
}

func TestContextEncryptDecryptStruct(t *testing.T) {
	div := Div()
	disp := NewClientTester(div)
	defer disp.Close()
	ctx := makeContext(div)

	expected := struct {
		Title string
		Value int
	}{
		Title: "hello",
		Value: 42,
	}

	item := expected
	item.Title = ""
	item.Value = 0

	crypted, err := ctx.Encrypt(expected)
	require.NoError(t, err)
	require.NotEmpty(t, crypted)

	err = ctx.Decrypt(crypted, &item)
	require.NoError(t, err)
	require.Equal(t, expected, item)
}

func TestContextEncryptDecryptString(t *testing.T) {
	div := Div()
	disp := NewClientTester(div)
	defer disp.Close()
	ctx := makeContext(div)

	expected := "hello"
	item := ""

	crypted, err := ctx.Encrypt(expected)
	require.NoError(t, err)
	require.NotEmpty(t, crypted)

	err = ctx.Decrypt(crypted, &item)
	require.NoError(t, err)
	require.Equal(t, expected, item)
}

func TestContextEncryptDecryptBytes(t *testing.T) {
	div := Div()
	disp := NewClientTester(div)
	defer disp.Close()
	ctx := makeContext(div)

	expected := []byte("hello")
	var item []byte

	crypted, err := ctx.Encrypt(expected)
	require.NoError(t, err)
	require.NotEmpty(t, crypted)

	err = ctx.Decrypt(crypted, &item)
	require.NoError(t, err)
	require.Equal(t, expected, item)
}

func TestContextHandle(t *testing.T) {
	foo := &foo{}
	client := NewClientTester(foo)
	defer client.Close()

	actionName := "/test/context/handle"
	action := Action{}
	ctx := makeContext(foo)

	ctx.Handle(actionName, func(ctx Context, a Action) {
		action = a
	})

	ctx.NewAction(actionName).
		Value(21).
		Tag("hello", "world").
		Post()

	client.Consume()
	require.Equal(t, actionName, action.Name)
	require.Equal(t, 21, action.Value)
	require.Equal(t, "world", action.Tags.Get("hello"))

}
