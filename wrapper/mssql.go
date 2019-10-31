package wrapper

// Init TracerWrapper with default ignoreSelectColumnsOption.
func NewMsSQLTracerWrapper() *TracerWrapper {
	return NewTracerWrapper("mssql")
}

// Init pure TracerWrapper with set options.
func NewMsSQLTracerWrapperWithOpts(options ...TracerOption) *TracerWrapper {
	return newTracerWrapper(newTracer("mssql", options...))
}
