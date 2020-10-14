package cog

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/kovetskiy/lorg"
	"github.com/reconquest/karma-go"
	"github.com/stretchr/testify/assert"
)

func TestFatal(t *testing.T) {
	test := assert.New(t)

	_ = test

	outBuffer := bytes.NewBuffer(nil)
	out := lorg.NewLog()
	out.SetOutput(outBuffer)
	out.SetFormat(lorg.NewFormat(`${level} %s`))

	logger := NewLogger(out)

	senderBuffer := bytes.NewBuffer(nil)

	logger.SetExiter(func(code int) {
		senderBuffer.WriteString(fmt.Sprintf("exit %d\n", code))
	})

	logger.SetSender(func(lvl lorg.Level, msg karma.Hierarchical) error {
		senderBuffer.WriteString(lvl.String() + " " + msg.String() + "\n")
		return nil
	})
	logger.Warning("warning")
	logger.Fatal("fatal")

	test.EqualValues(
		"WARNING warning\n"+
			"FATAL fatal\n"+
			"exit 1\n",
		senderBuffer.String(),
	)
}
