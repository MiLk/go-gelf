package message

import (
	"bytes"
	"encoding/json"
	"log/syslog"
	"time"
)

// Message represents the contents of the GELF message.
type Message struct {
	Version  string                 `json:"version"`
	Host     string                 `json:"host"`
	Short    string                 `json:"short_message"`
	Full     string                 `json:"full_message,omitempty"`
	TimeUnix float64                `json:"timestamp"`
	Level    syslog.Priority        `json:"level,omitempty"`
	Extra    map[string]interface{} `json:"-"`
	RawExtra json.RawMessage        `json:"-"`
}

func (m *Message) UnmarshalJSON(data []byte) error {
	i := make(map[string]interface{}, 16)
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	for k, v := range i {
		if k[0] == '_' {
			if m.Extra == nil {
				m.Extra = make(map[string]interface{}, 1)
			}
			m.Extra[k] = v
			continue
		}
		switch k {
		case "version":
			m.Version = v.(string)
		case "host":
			m.Host = v.(string)
		case "short_message":
			m.Short = v.(string)
		case "full_message":
			m.Full = v.(string)
		case "timestamp":
			m.TimeUnix = v.(float64)
		case "level":
			m.Level = syslog.Priority(v.(float64))
		}
	}
	return nil
}

func (m *Message) Bytes() ([]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	hasExtra := len(m.Extra) > 0
	hasRawExtra := len(m.RawExtra) > 0

	if !hasExtra && !hasRawExtra {
		return b, nil
	}

	parts := [][]byte{b[:len(b)-1]}

	if hasExtra {
		eb, err := json.Marshal(m.Extra)
		if err != nil {
			return nil, err
		}

		parts = append(parts, []byte(","), eb[1:len(eb)-1])
	}

	if hasRawExtra {
		parts = append(parts, []byte(","), m.RawExtra[1:len(m.RawExtra)-1])
	}

	parts = append(parts, []byte("}"))

	totalLength := 0
	for _, p := range parts {
		totalLength += len(p)
	}

	res := make([]byte, totalLength)

	var i int
	for _, p := range parts {
		i += copy(res[i:], p)
	}

	return res, nil
}

func New(p []byte, host string, file string, line int) *Message {
	// remove trailing and leading whitespace
	p = bytes.TrimSpace(p)

	// If there are newlines in the message, use the first line
	// for the short message and set the full message to the
	// original input.  If the input has no newlines, stick the
	// whole thing in Short.
	short := p
	full := []byte("")
	if i := bytes.IndexRune(p, '\n'); i > 0 {
		short = p[:i]
		full = p
	}

	return &Message{
		Version:  "1.1",
		Host:     host,
		Short:    string(short),
		Full:     string(full),
		TimeUnix: float64(time.Now().Unix()),
		Level:    6, // info
		Extra: map[string]interface{}{
			"_file": file,
			"_line": line,
		},
	}
}
