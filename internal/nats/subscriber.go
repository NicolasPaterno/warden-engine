package nats

import (
	"context"
	"log/slog"
	"time"

	sensorv1 "github.com/NicolasPaterno/warden-proto/gen/go/warden/sensor/v1"
	natsgo "github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/proto"
)

const subject = "warden.sensors.v1.>"

type Subscriber struct {
	conn *natsgo.Conn
}

func NewSubscriber(url string) (*Subscriber, error) {
	conn, err := natsgo.Connect(
		url,
		natsgo.ReconnectHandler(
			func(conn *natsgo.Conn) {
				slog.Info("reconnected to NATS", "url", conn.ConnectedUrl())
			}),
		natsgo.MaxReconnects(10),
		natsgo.ReconnectWait(5*time.Second),
		natsgo.RetryOnFailedConnect(true),
		natsgo.ConnectHandler(func(conn *natsgo.Conn) {
			slog.Info("connected to NATS", "url", conn.ConnectedUrl())
		}),
		natsgo.DisconnectErrHandler(func(conn *natsgo.Conn, err error) {
			slog.Warn("disconnected from NATS", "error", err)
		}))
	if err != nil {
		return nil, err
	}
	return &Subscriber{conn: conn}, nil
}

func (s *Subscriber) Subscribe(ctx context.Context, handler func(ctx context.Context, reading *sensorv1.SensorReading)) error {
	_, err := s.conn.Subscribe(subject, func(msg *natsgo.Msg) {
		msgCtx := otel.GetTextMapPropagator().Extract(ctx, natsHeaderCarrier{msg.Header})
		var reading sensorv1.SensorReading
		if err := proto.Unmarshal(msg.Data, &reading); err != nil {
			slog.Error("Failed to unmarshal reading", "error", err)
			return
		}
		handler(msgCtx, &reading)
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Subscriber) Close() error {
	return s.conn.Drain()
}

type natsHeaderCarrier struct {
	header natsgo.Header
}

func (c natsHeaderCarrier) Get(key string) string {
	return c.header.Get(key)
}

func (c natsHeaderCarrier) Set(key, value string) {
	c.header.Set(key, value)
}

func (c natsHeaderCarrier) Keys() []string {
	result := make([]string, 0, len(c.header))
	for k := range c.header {
		result = append(result, k)
	}
	return result
}

func (s *Subscriber) IsConnected() bool {
	return s.conn.IsConnected()
}
