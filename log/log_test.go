package log

import (
	"testing"
)

func TestSetLevel(t *testing.T) {
	t.Run("infoLevel", func(t *testing.T) {
		SetLevel(InfoLevel)
		Info("this is a info")
		Error("this is a error")
	})
	t.Run("errorLevel", func(t *testing.T) {
		SetLevel(ErrorLevel)
		Info("this is a info")
		Error("this is a error")
	})
	t.Run("disabled", func(t *testing.T) {
		SetLevel(Disabled)
		Info("this is a info")
		Error("this is a error")
	})
}
