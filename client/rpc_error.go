package client

import (
	"encoding/json"
	"fmt"
	"io"
)

// RPCError is satisfied by every error surfaced from a specular call: built-in
// errors, user-defined errors (which are structs implementing the error
// interface) and the generic UnknownRPCError produced when the client cannot
// resolve the concrete error type.
type RPCError interface {
	error
	Struct
}

// UnknownRPCError is surfaced when the server returns an error whose type the
// client was not generated with knowledge of (for example a newer server
// introduced an error type after the client was generated).
//
// It still satisfies the error interface so callers handle it like any other
// error at the time exposing the raw payload so forward-compatible code can read the
// fields it knows.
type UnknownRPCError struct {
	// Path is the fully-qualified type path (ns.mod.type) advertised by the server.
	Path StructPath
	// Message is the human-readable message, when discoverable in the payload.
	Message string
	// Code is the machine-readable code, when discoverable in the payload.
	Code string
	// Payload is the raw, undecoded error body.
	Payload json.RawMessage
}

// StructPath implements the Struct interface.
func (e *UnknownRPCError) StructPath() StructPath {
	return e.Path
}

// Error implements the error interface.
func (e *UnknownRPCError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return "unknown rpc error " + e.Path.String()
}

// newUnknownRPCError reads the error body and builds a generic error for an
// unresolved type. field extraction and payload is both preserved best-effort
func newUnknownRPCError(path StructPath, body io.Reader) (*UnknownRPCError, error) {
	raw, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("failed to read unknown error body, %w", err)
	}
	ue := &UnknownRPCError{
		Path:    path,
		Payload: json.RawMessage(raw),
	}
	var fields struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	}
	_ = json.Unmarshal(raw, &fields)
	ue.Message = fields.Message
	ue.Code = fields.Code
	return ue, nil
}
