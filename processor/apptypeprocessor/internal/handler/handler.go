package handler

import (
	"go.opentelemetry.io/collector/pdata/pcommon"
	conventions "go.opentelemetry.io/collector/semconv/v1.22.0"
)

const (
	RocketmqConnectionAttrName = "messaging.rocketmq.broker_address"
	ActivemqConnectionAttrName = "messaging.activemq.broker_address"
)

type handler interface {
	GetHostPort(attrs *pcommon.Map) (host string, port string)
}

var handlers map[string]handler

func InitHandlers() {
	if handlers == nil {
		handlers = make(map[string]handler)

		handlers["rocketmq"] = &DefaultComponentHandler{UrlKey: RocketmqConnectionAttrName}
		handlers["rabbitmq"] = &RabbitmqHandler{}
		handlers["activemq"] = &ActivemqHandler{}
		handlers["memcached"] = &DefaultComponentHandler{UrlKey: conventions.AttributeDBConnectionString}

		handlers["default"] = &DefaultDBHandler{}
	}
}

func GetHandler(name string) handler {
	h, ok := handlers[name]
	if !ok {
		return handlers["default"]
	} else {
		return h
	}
}
