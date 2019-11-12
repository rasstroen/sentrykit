package sentrykit

import (
	"errors"
	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func Test_getValue(t *testing.T) {
	type args struct {
		key     string
		keyvals []interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "key exists",
			args: args{
				key:     "test",
				keyvals: []interface{}{"test", "testValue"},
			},
			want: "testValue",
		},
		{
			name: "second key exists",
			args: args{
				key:     "test2",
				keyvals: []interface{}{"test", "testValue", "test2", "testValue2"},
			},
			want: "testValue2",
		},
		{
			name: "value does not exists",
			args: args{
				key:     "test2",
				keyvals: []interface{}{"test", "testValue", "test2"},
			},
			want: nil,
		},
		{
			name: "key does not exist mean",
			args: args{
				key:     "test3",
				keyvals: []interface{}{"test", "testValue", "test2", "testValue2"},
			},
			want: nil,
		},
		{
			name: "key does not exist odd",
			args: args{
				key:     "test3",
				keyvals: []interface{}{"test", "testValue", "test2"},
			},
			want: nil,
		},
		{
			name: "empty",
			args: args{
				key:     "test3",
				keyvals: []interface{}{},
			},
			want: nil,
		},
		{
			name: "tail problem",
			args: args{
				key:     "test",
				keyvals: []interface{}{"key1", "test", "key2", "test"},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getValue(tt.args.key, tt.args.keyvals); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

type sentryException struct {
	exception error
	hint      *sentry.EventHint
	scope     sentry.EventModifier
}

type sentryMessage struct {
	message string
	hint    *sentry.EventHint
	scope   sentry.EventModifier
}
type sentryClientMock struct {
	exceptions []sentryException
	messages   []sentryMessage
}

func newEventId(id string) *sentry.EventID {
	eventId := sentry.EventID(id)
	return &eventId
}

func (s *sentryClientMock) CaptureException(exception error, hint *sentry.EventHint, scope sentry.EventModifier) *sentry.EventID {
	s.exceptions = append(s.exceptions, sentryException{exception: exception, hint: hint, scope: scope})
	return newEventId("id")
}

func (s *sentryClientMock) CaptureMessage(message string, hint *sentry.EventHint, scope sentry.EventModifier) *sentry.EventID {
	s.messages = append(s.messages, sentryMessage{message: message, hint: hint, scope: scope})
	return newEventId("id")
}

func TestSentryLogger_Log(t *testing.T) {
	tests := []struct {
		name           string
		keyvals        []interface{}
		wantErr        bool
		wantExceptions []sentryException
		wantMessages   []sentryMessage
	}{
		{
			name:           "empty",
			keyvals:        []interface{}{},
			wantErr:        false,
			wantExceptions: []sentryException{},
			wantMessages: []sentryMessage{
				{
					message: "",
					hint:    &sentry.EventHint{},
					scope:   scopeWithExtra([]interface{}{}),
				},
			},
		},
		{
			name:    "error",
			keyvals: []interface{}{"err", errors.New("some error")},
			wantErr: false,
			wantExceptions: []sentryException{
				{
					exception: errors.New("some error"),
					hint:      &sentry.EventHint{},
					scope:     scopeWithExtra([]interface{}{"err", errors.New("some error")}),
				},
			},
			wantMessages: []sentryMessage{},
		}, {
			name:           "msg",
			keyvals:        []interface{}{"msg", "test message"},
			wantErr:        false,
			wantExceptions: []sentryException{},
			wantMessages: []sentryMessage{
				{
					message: "test message",
					hint:    &sentry.EventHint{},
					scope:   scopeWithExtra([]interface{}{"msg", "test message"}),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sentryClientMock := sentryClientMock{
				exceptions: []sentryException{},
				messages:   []sentryMessage{},
			}
			s := NewSentryLogger(&sentryClientMock)

			if err := s.Log(tt.keyvals...); (err != nil) != tt.wantErr {
				t.Errorf("Log() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.wantExceptions, sentryClientMock.exceptions)
			assert.Equal(t, tt.wantMessages, sentryClientMock.messages)
		})
	}
}

func Test_convertToKeyValueMap(t *testing.T) {
	tests := []struct {
		name    string
		keyvals []interface{}
		want    map[string]interface{}
	}{
		{
			name:    "simple",
			keyvals: []interface{}{"key", "value", "key1", "value1", "count", 1},
			want:    map[string]interface{}{"key": "value", "key1": "value1", "count": 1},
		},
		{
			name:    "odd",
			keyvals: []interface{}{"key", "value", "key1"},
			want:    map[string]interface{}{"key": "value"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertToKeyValueMap(tt.keyvals); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertToKeyValueMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
