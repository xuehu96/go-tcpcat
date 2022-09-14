package snowflake

import (
	"testing"
)

func TestSnowflake(t *testing.T) {
	Init("2020-01-01", 1)
	t.Logf("%s", GetSn())
	t.Logf("%d", GetId())
}
