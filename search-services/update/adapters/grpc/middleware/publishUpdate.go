package middleware

import (
	"context"
	"log/slog"
	"slices"

	"google.golang.org/grpc"
)

type UpdatePublisher interface {
	PublishUpdateEvent(updateEventName string) error
}

func PublishUpdate(
	publisher UpdatePublisher,
	updateEventName string,
	log *slog.Logger,
	methods ...string,
) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		resp, err = handler(ctx, req)
		if err != nil {
			return resp, err
		}

		if !slices.Contains(methods, info.FullMethod) {
			return resp, err
		}

		log.Debug("publish method", "method", info.FullMethod)

		err = publisher.PublishUpdateEvent(updateEventName)
		if err != nil {
			// не стал возвращать ошибку на клиент, так как тогда не понятно,
			// было ли совершено само действие (update, например)
			log.Error("failed to publish update event", "error", err)
		}

		return resp, nil
	}
}
