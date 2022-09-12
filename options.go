package zaptool

type loggingOptions struct {
	includeTiming    bool
	includeTimestamp bool
}

type loggingOptionsFunc func(o *loggingOptions)

// LoggingOptionTiming defines if the logging should contain a `http.request_time` field.
//
//nolint:revive // deliberately not-exported function type.
func LoggingOptionTiming(state bool) loggingOptionsFunc {
	return func(o *loggingOptions) {
		o.includeTiming = state
	}
}

// LoggingOptionTimestamp defines if the logging should contain a `http.timestamp` field.
//
//nolint:revive // deliberately not-exported function type.
func LoggingOptionTimestamp(state bool) loggingOptionsFunc {
	return func(o *loggingOptions) {
		o.includeTimestamp = state
	}
}
