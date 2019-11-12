package sentrykit

import (
	"github.com/getsentry/sentry-go"
)

type SentryClient interface {
	CaptureException(exception error, hint *sentry.EventHint, scope sentry.EventModifier) *sentry.EventID
	CaptureMessage(message string, hint *sentry.EventHint, scope sentry.EventModifier) *sentry.EventID
}

type SentryLogger struct {
	client SentryClient
}

// Sentry Logger captures exceptions if you provide "err" field
// other wise it captures simple messages
func (s *SentryLogger) Log(keyvals ...interface{}) error {
	scope := scopeWithExtra(keyvals)
	hint := &sentry.EventHint{}

	err, ok := getValue("err", keyvals).(error)
	if ok {
		s.client.CaptureException(
			err,
			hint,
			scope,
		)
	} else {
		msg, ok := getValue("msg", keyvals).(string)
		if !ok {
			msg = ""
		}
		s.client.CaptureMessage(
			msg,
			hint,
			scopeWithExtra(keyvals),
		)
	}

	return nil
}

func NewSentryLogger(client SentryClient) *SentryLogger {
	return &SentryLogger{client: client}
}

func scopeWithExtra(keyvals []interface{}) *sentry.Scope {
	scope := sentry.NewScope()
	scope.SetExtras(convertToKeyValueMap(keyvals))

	return scope
}

func convertToKeyValueMap(keyvals []interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(keyvals)/2)
	for i := 0; i < len(keyvals)-1; i += 2 {
		currentKey, ok := keyvals[i].(string)
		if ok {
			result[currentKey] = keyvals[i+1]
		}
	}
	return result
}

func getValue(key string, keyvals []interface{}) interface{} {
	for i := 0; i < len(keyvals)-1; i += 2 {
		currentKey, ok := keyvals[i].(string)
		if !ok {
			continue
		}

		if currentKey == key {
			return keyvals[i+1]
		}
	}

	return nil
}
