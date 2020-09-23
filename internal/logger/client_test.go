package logger

import (
	"encoding/json"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	logger := NewForTest()
	message := "an example warning message"
	now := time.Now()
	logger.WithTime(now).Warning(message)

	out, ok := logger.Out.(TestOutput)
	if !ok {
		t.Fatal("Unable to transform type Logger.Out to TestOutput")
	}

	if len(*out.Written) == 0 {
		t.Fatal("Log written buffer contains 0 bytes")
	}

	actual := out.Lines()[0]

	j, _ := json.Marshal(map[string]string{
		"level": "warning",
		"msg":   message,
		"time":  now.Format(time.RFC3339),
	})
	expected := string(j) // {"level":"warning","msg":"an example warning message","time":"2020-09-06T04:46:13+03:00"}

	if actual != expected {
		t.Errorf("Logged message not equal: expected = %s, actual = %s", expected, actual)
	}
}
