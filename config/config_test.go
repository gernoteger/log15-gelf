package config

import (
	"fmt"
	"testing"

	"net"

	"github.com/davecgh/go-spew/spew"
	"github.com/gernoteger/log15-config"
	"github.com/gernoteger/log15-gelf"
	"github.com/inconshreveable/log15"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func testConfigLogger(conf string) (log15.Logger, error) {
	configMap := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(conf), &configMap)
	if err != nil {
		return nil, fmt.Errorf("yaml umnarshall failed: %v", err)
	}
	return config.Logger(configMap)
}

// givenHostAvailable skips if host not available
func givenHostAvailable(host string, t *testing.T) {
	ips, err := net.LookupIP(host)
	if err != nil || len(ips) == 0 {
		t.Skipf("can't resolve host '%v'", host)
	}
}

func TestGelfConfig(t *testing.T) {
	t.Parallel()

	require := require.New(t)
	assert := assert.New(t)

	r, err := gelf.NewReader("127.0.0.1:0")
	if err != nil {
		t.Fatalf("NewReader faild: %v", err)
	}

	r.Addr()
	config := fmt.Sprintf(`
  level: INFO
  handlers:
    - kind: gelf
      address: "%v"
`, r.Addr())

	fmt.Println(r.Addr())
	//givenHostAvailable("logger", t)

	l, err := testConfigLogger(config)
	require.Nil(err)

	l.Info("Hello, gelf!", "mark", 1)

	msg, err := r.ReadMessage()

	spew.Dump(msg)
	assert.Equal("Hello, gelf!", msg.Short)
	assert.EqualValues(1, msg.Extra["_mark"])
}
