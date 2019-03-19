package log4go

type (
	OutputMode string

	LogExtInfo struct {
		modulename string
		outputmode OutputMode
	}
)

const (
	OutputModeDefault OutputMode = ""
	OutputModeJson    OutputMode = "json"
)

func (w *FileLogWriter) SetModuleName(modulename string) (*FileLogWriter) {
	w.ext_info.modulename = modulename
	return w
}

func (w *FileLogWriter) SetOutputMode(outputmode OutputMode) (*FileLogWriter){
	w.ext_info.outputmode = outputmode
	return w
}

func (w *FileLogWriter) formatLogRecord(format string, rec *LogRecord) string {
	if w.ext_info.outputmode == OutputModeJson {
		return FormatLogRecordJson(format, w.ext_info.modulename, rec)
	} else {
		return FormatLogRecord(format, rec)
	}
}
