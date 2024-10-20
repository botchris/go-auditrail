package auditrail

import (
	"time"

	"github.com/google/uuid"
)

// Entry represents an audit log event.
//
// This struct is not safe for concurrent write access.
type Entry struct {
	// IdempotencyID is used to uniquely identify the log entry, it used as
	// deduplication key when streaming logs to a log aggregator.
	IdempotencyID string `json:"idempotency_id,omitempty"`

	// Actor is who is making the action, it could be a username, an identifier,
	// etc.
	Actor string `json:"actor"`

	// Action is the action being performed, for example: "order_create",
	// "password_changed", etc.
	Action string `json:"action"`

	// Module is the module that the action is being performed on, it could be a
	// service name, a package name, etc. e.g. "users", "orders", etc.
	Module string `json:"module"`

	// CorrelationID is used to correlate multiple log entries that are related
	// to the same action.
	CorrelationID string `json:"correlation_id,omitempty"`

	// CausationID is used to track the original action that caused the current
	// action to be performed.
	CausationID string `json:"causation_id,omitempty"`

	// AuthMethod is the method that was used to authenticate the actor.
	AuthMethod string `json:"auth_method,omitempty"`

	// Details is used to attach any additional information to the log entry.
	Details map[string]interface{} `json:"details,omitempty"`

	// OccurredAt when the action was performed.
	OccurredAt time.Time `json:"occurred_at"`
}

// NewEntry creates a new log entry with the given actor, action, and module.
//
// By default, the following fields are set:
//
//   - The time of the event is set to the current time, use
//     [Entry.WithOccurredAt] to override it to a different value.
//   - The IdempotencyID is set to randomly generated value, use
//     [Entry.WithIdempotency] to override it to a different value.
func NewEntry(actor, action, module string) *Entry {
	return &Entry{
		IdempotencyID: uuid.NewString(),
		Actor:         actor,
		Action:        action,
		Module:        module,
		OccurredAt:    time.Now(),
	}
}

// WithIdempotency sets the idempotency ID of the event.
func (e *Entry) WithIdempotency(idempotencyID string) *Entry {
	e.IdempotencyID = idempotencyID

	return e
}

// WithCorrelation sets the correlation ID of the event.
func (e *Entry) WithCorrelation(correlationID string) *Entry {
	e.CorrelationID = correlationID

	return e
}

// WithCausation sets the causation ID of the event.
func (e *Entry) WithCausation(causationID string) *Entry {
	e.CausationID = causationID

	return e
}

// WithAuthMethod sets the authentication method of the event.
func (e *Entry) WithAuthMethod(method string) *Entry {
	e.AuthMethod = method

	return e
}

// AppendDetails adds a key-value pair to the details of the event.
func (e *Entry) AppendDetails(key string, value interface{}) *Entry {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}

	e.Details[key] = value

	return e
}

// WithOccurredAt overrides the time of the event with the given time.
func (e *Entry) WithOccurredAt(when time.Time) *Entry {
	e.OccurredAt = when

	return e
}
