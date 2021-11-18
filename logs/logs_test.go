package logs

import "testing"

func createLogs() *Logger {
	return NewLogger("test_logs")
}

func TestError(t *testing.T) {
	logeer := createLogs()
	data := make(map[string]interface{})
	data["class"] = "error"
	data["msg"] = "test"
	logeer.Error(data, "this error")
}

func TestDebug(t *testing.T) {
	Logger := createLogs()
	data := make(map[string]interface{})
	data["class"] = "error"
	data["msg"] = "test"
	Logger.Debug(data, "this debug")
}

func TestInfo(t *testing.T) {
	Logger := createLogs()
	data := make(map[string]interface{})
	data["class"] = "error"
	data["msg"] = "test"
	Logger.Info(data, "this Info")
}

func TestWarn(t *testing.T) {
	Logger := createLogs()
	data := make(map[string]interface{})
	data["class"] = "error"
	data["msg"] = "test"
	Logger.Warn(data, "this Warn")
}
