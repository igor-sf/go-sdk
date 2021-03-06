package client

import (
	"context"

	pb "github.com/dapr/go-sdk/dapr/proto/runtime/v1"
	"github.com/pkg/errors"
)

// BindingInvocation represents binding invocation request
type BindingInvocation struct {
	// Name is name of binding to invoke.
	Name string
	// Operation is the name of the operation type for the binding to invoke
	Operation string
	// Data is the input bindings sent
	Data []byte
	// Metadata is the input binding metadata
	Metadata map[string]string
}

// BindingEvent represents the binding event handler input
type BindingEvent struct {
	// Data is the input bindings sent
	Data []byte
	// Metadata is the input binding metadata
	Metadata map[string]string
}

// InvokeBinding invokes specific operation on the configured Dapr binding.
// This method covers input, output, and bi-directional bindings.
func (c *GRPCClient) InvokeBinding(ctx context.Context, in *BindingInvocation) (out *BindingEvent, err error) {
	if in == nil {
		return nil, errors.New("binding invocation required")
	}
	if in.Name == "" {
		return nil, errors.New("binding invocation name required")
	}
	if in.Operation == "" {
		return nil, errors.New("binding invocation operation required")
	}

	req := &pb.InvokeBindingRequest{
		Name:      in.Name,
		Operation: in.Operation,
		Data:      in.Data,
		Metadata:  in.Metadata,
	}

	resp, err := c.protoClient.InvokeBinding(authContext(ctx), req)
	if err != nil {
		return nil, errors.Wrapf(err, "error invoking binding %s/%s", in.Name, in.Operation)
	}

	out = &BindingEvent{}

	if resp != nil {
		out.Data = resp.Data
		out.Metadata = resp.Metadata
	}

	return
}

// InvokeOutputBinding invokes configured Dapr binding with data (allows nil).InvokeOutputBinding
// This method differs from InvokeBinding in that it doesn't expect any content being returned from the invoked method.
func (c *GRPCClient) InvokeOutputBinding(ctx context.Context, in *BindingInvocation) error {
	if _, err := c.InvokeBinding(ctx, in); err != nil {
		return errors.Wrap(err, "error invoking output binding")
	}
	return nil
}
