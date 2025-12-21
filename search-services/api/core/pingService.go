package core

import "context"

const (
	ServiceStatusOk          = "ok"
	ServiceStatusUnavailable = "unavailable"
)

func PingServices(ctx context.Context, pingers map[string]Pinger) map[string]string {
	services := make(map[string]string)
	for name, service := range pingers {
		err := service.Ping(ctx)
		if err != nil {
			services[name] = ServiceStatusUnavailable
		} else {
			services[name] = ServiceStatusOk
		}
	}

	return services
}
