package messaging

import "context"

type Producer interface {
	Send(ctx context.Context, msg any) error
}
