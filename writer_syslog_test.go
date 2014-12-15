package lsyslog

import (
	"testing"

	gc "gopkg.in/check.v1"

	"github.com/juju/loggo"
)

func Test(t *testing.T) {
	gc.TestingT(t)
}

type loggerSuite struct{}

var _ = gc.Suite(&loggerSuite{})

func (*loggerSuite) SetUpTest(c *gc.C) {
	loggo.ResetLoggers()
	loggo.ResetWriters()
}

// note this is really just a smoke test at the moment.
func (*loggerSuite) TestSyslogWriter(c *gc.C) {
	defaultWriter, level, err := loggo.RemoveWriter("default")
	c.Assert(err, gc.IsNil)
	c.Assert(level, gc.Equals, loggo.TRACE)
	c.Assert(defaultWriter, gc.NotNil)

	err = loggo.RegisterWriter("default", NewDefaultSyslogWriter(loggo.INFO, "", "LOCAL7"), loggo.INFO)
	c.Assert(err, gc.IsNil)

	loggo.GetLogger("").Infof("some test %s", "hello")
}
