package test

import (
	"fmt"
	"math"
	"os"
	"testing"

	config "github.com/carlosvin/go-config-example/internal"
	"github.com/stretchr/testify/assert"
)

type TestConfig struct {
	A   string  `cfg_name:"a"`
	B   int     `cfg_name:"b" cfg_desc:"B is a flag for someting"`
	B64 int64   `cfg_name:"b64" cfg_desc:"B 64 value for something"`
	C   float64 `cfg_name:"c" cfg_desc:"C is a config value for something"`
}

func TestEnv(t *testing.T) {
	tc := &TestConfig{A: "default", B: 123, C: 3.1415926}
	cfg := config.New(tc, []string{})

	os.Setenv("B", "999")
	cfg.Load()
	assert.Equal(t, "default", tc.A)
	assert.Equal(t, 999, tc.B)
	assert.Equal(t, int64(0), tc.B64)
	assert.Equal(t, 3.1415926, tc.C)

	os.Setenv("A", "new value")
	os.Setenv("B64", fmt.Sprintf("%d", math.MaxInt64))
	os.Setenv("C", "0.99999")
	cfg.Load()
	assert.Equal(t, "new value", tc.A)
	assert.Equal(t, 999, tc.B)
	assert.Equal(t, int64(math.MaxInt64), tc.B64)
	assert.Equal(t, 0.99999, tc.C)
	cleanupEnv()
}

func TestFlags(t *testing.T) {
	args := []string{
		"test-flags",
		"--a=imaflag",
		fmt.Sprintf("--b64=%d", math.MinInt64),
	}

	tc := &TestConfig{A: "default", B: 123, B64: 1, C: 1.1111}
	cfg := config.New(tc, args)
	cfg.Load()

	assert.Equal(t, "imaflag", tc.A)
	assert.Equal(t, 123, tc.B)
	assert.Equal(t, int64(math.MinInt64), tc.B64)
	assert.Equal(t, 1.1111, tc.C)

}

func cleanupEnv() {
	os.Setenv("A", "")
	os.Setenv("B64", "")
	os.Setenv("B", "")
	os.Setenv("C", "")
}
