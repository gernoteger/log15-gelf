package gelf

import (
	"testing"
	"time"

	"github.com/go-stack/stack"
	"github.com/inconshreveable/log15"
)

func TestCtxToMap(t *testing.T) {

	loc, err := time.LoadLocation("Europe/Vienna")
	if err != nil {
		t.Fatalf("can't load Timezone: %v", err)
	}
	logTime := time.Date(2016, 11, 23, 13, 01, 02, 123100*1e3, loc)

	expected := map[string]interface{}{
		"_msg":    "a message",
		"_foo":    "baz",
		"_number": 1,
		"_t":      logTime,
	}
	ctx := []interface{}{"msg", "a message", "foo", "bar", "foo", "baz", "number", 1, "t", logTime}

	cm := ctxToMap(ctx)

	for k, v := range expected {
		if cm[k] != v {
			t.Fatalf("%v: expected: '%v', got: %v", k, v, cm[k])
		}
	}
}

const SyslogInfoLevel = 6

func TestGelfHandler(t *testing.T) {
	t.Parallel()

	r, err := NewReader("127.0.0.1:0")
	if err != nil {
		t.Fatalf("NewReader faild: %v", err)
	}

	loc, err := time.LoadLocation("Europe/Vienna")
	if err != nil {
		t.Fatalf("can't load Timezone: %v", err)
	}
	logTime := time.Date(2016, 11, 23, 13, 01, 02, 123100*1e3, loc)

	msgData := "test message\nsecond line"
	rec := log15.Record{
		Time: logTime, //TODO: set fixed!!
		Lvl:  log15.LvlInfo,
		Msg:  msgData,
		Ctx:  []interface{}{"foo", "bar", "withField", "1", "foo", "baz"}, // no fields yet
		Call: stack.Caller(0),
	}

	h := Must.GelfHandler(r.Addr())

	h.Log(&rec)

	msg, err := r.ReadMessage()

	if err != nil {
		t.Fatalf("Couldn't read Message: %v", err)
	}
	if msg.Short != "test message" {
		t.Fatalf("msg.Short expected: 'text message', got: %v", msg.Short)
	}
	if msg.Full != msgData {
		t.Fatalf("msg.Full expected: '%v', got: %v", msgData, msg.Full)
	}
	if msg.Level != SyslogInfoLevel {
		t.Fatalf("msg.Level expected: '%v', got: %v", SyslogInfoLevel, msg.Level)
	}
	if len(msg.Extra) != 2 {
		t.Fatalf("msg.Extra length expected: '%v', got: %v", 2, len(msg.Extra))
	}
	if msg.File != "gelf_test.go" {
		t.Fatalf("msg.File expected: '%v', got: %v", "gelf_test.go", msg.File)
	}

	// no tests for line; this would be too unstable..
	extra := map[string]string{"foo": "baz", "withField": "1"}

	for k, v := range extra {
		// extra fields are prefixed with "_"
		val, ok := msg.Extra["_"+k].(string)
		if !ok {
			t.Fatalf("no key foundfor %v", k)
		}
		if val != v {
			t.Fatalf("extra[%v] expected: '%v', got: '%v'", k, v, val)
		}
	}

	// checking time...
	s := int64(msg.TimeUnix)
	ns := int64((msg.TimeUnix - float64(s)) * 1e9)
	mt := time.Unix(s, ns)

	//fmt.Printf("t0=%v time=%v t=%v", msg.TimeUnix, mt, logTime); fmt.Println()
	diff := logTime.Sub(mt)
	//fmt.Printf("diff=%v", diff)
	//fmt.Println()
	if !within(diff, time.Millisecond) {
		t.Fatalf("difference too big: %v", diff)
	}
	//assert.WithinDuration(logTime, mt, time.Millisecond, "time from log") // we have millisecond precision

}

func within(dt time.Duration, delta time.Duration) bool {
	if dt < -delta || dt > delta {
		return false
	}
	return true
}
