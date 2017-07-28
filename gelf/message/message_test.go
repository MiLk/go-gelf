package message

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMessage_Bytes(t *testing.T) {
	m := Message{
		Version:  "1.1",
		Host:     "test-host",
		Short:    "short message",
		Full:     "full message",
		TimeUnix: float64(time.Date(2017, 7, 7, 7, 7, 7, 7, time.UTC).Unix()),
		Level:    6, // info
		Extra:    map[string]interface{}{"_file": "1234", "_line": "3456"},
	}
	b, err := m.Bytes()
	assert.Nil(t, err)
	expected := `{"version":"1.1","host":"test-host","short_message":"short message","full_message":"full message","timestamp":1499411227,"level":6,"_file":"1234","_line":"3456"}`
	assert.EqualValues(t, expected, string(b))
}

func BenchmarkMessage_Bytes(b *testing.B) {
	host, err := os.Hostname()
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := Message{
			Version:  "1.1",
			Host:     host,
			Short:    "short message",
			Full:     "full message",
			TimeUnix: float64(time.Now().Unix()),
			Level:    6, // info
			Extra:    map[string]interface{}{"_file": "1234", "_line": "3456"},
		}
		m.Bytes()
	}
}
