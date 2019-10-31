package wrapper

// Init TracerWrapper with default ignoreSelectColumnsOption.
func NewMySQLTracerWrapper() *TracerWrapper {
	return NewTracerWrapper("mysql")
}

// Init pure TracerWrapper with set options.
func NewMySQLTracerWrapperWithOpts(options ...TracerOption) *TracerWrapper {
	return newTracerWrapper(newTracer("mysql", options...))
}
