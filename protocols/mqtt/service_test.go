package mqtt

import (
	"encoding/json"
	"testing"
)

// fakeMessage implements mqtt.Message for testing purposes
type fakeMessage struct {
	topic   string
	payload []byte
}

func (f *fakeMessage) Duplicate() bool   { return false }
func (f *fakeMessage) Qos() byte         { return 0 }
func (f *fakeMessage) Retained() bool    { return false }
func (f *fakeMessage) Topic() string     { return f.topic }
func (f *fakeMessage) MessageID() uint16 { return 0 }
func (f *fakeMessage) Payload() []byte   { return f.payload }

func TestPublishHandlerMarshals(t *testing.T) {
	called := make(chan bool, 1)

	handler := func(b []byte) {
		var msg struct {
			Topic   string `json:"topic"`
			Payload string `json:"payload"`
		}
		if err := json.Unmarshal(b, &msg); err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}
		if msg.Topic != "test/topic" || msg.Payload != "hello" {
			t.Fatalf("unexpected message content: %+v", msg)
		}
		called <- true
	}

	// simulate the default publish handler behaviour used in startMqttSubscriber
	m := &fakeMessage{topic: "test/topic", payload: []byte("hello")}
	// construct json like the handler would
	data := struct {
		Topic   string `json:"topic"`
		Payload string `json:"payload"`
	}{
		Topic:   m.Topic(),
		Payload: string(m.Payload()),
	}
	b, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	// call handler
	handler(b)

	select {
	case <-called:
		// ok
	default:
		t.Fatalf("handler was not called")
	}
}
