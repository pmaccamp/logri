package logri_test

import (
	"github.com/Sirupsen/logrus"
	. "github.com/iancmcc/logri"
	. "gopkg.in/check.v1"
)

func (s *LogriSuite) TestSetLoggerLevel(c *C) {
	// Set the level to info and log below that
	s.logger.SetLevel(logrus.InfoLevel, true)
	s.AssertNotLogs(s.logger.Debug, "debug msg 1")
	// Now set the level to debug and log at debug
	s.logger.SetLevel(logrus.DebugLevel, true)
	s.AssertLogs(s.logger.Debug, "debug msg 2")
}

func (s *LogriSuite) TestUnsetLoggerLevel(c *C) {
	err := s.logger.SetLevel(NilLevel, true)
	c.Assert(err.Error(), Equals, ErrInvalidRootLevel.Error())

	alogger := s.logger.GetChild("a")
	alogger.SetLevel(logrus.ErrorLevel, true)

	s.AssertLogs(alogger.Error, "error message")

	err = alogger.SetLevel(NilLevel, true)
	c.Assert(err, IsNil)

	s.AssertNotLogs(alogger.Debug, "debug msg")
	s.AssertLogs(alogger.Info, "info msg")
}

func (s *LogriSuite) TestGetChildLogger(c *C) {
	alogger := s.logger.GetChild("a")
	c.Assert(alogger.Name, Equals, "a")

	blogger := alogger.GetChild("b")
	c.Assert(blogger.Name, Equals, "a.b")

	blogger2 := s.logger.GetChild("a.b")
	c.Assert(blogger2, Equals, blogger)

	clogger := s.logger.GetChild("a.b.c")
	c.Assert(clogger.Name, Equals, "a.b.c")

	clogger2 := alogger.GetChild("a.b.c")
	c.Assert(clogger2, Equals, clogger)

	clogger3 := blogger.GetChild("a.b.c")
	c.Assert(clogger3, Equals, clogger)

	clogger4 := blogger.GetChild("c")
	c.Assert(clogger4, Equals, clogger)

	clogger5 := alogger.GetChild("b.c")
	c.Assert(clogger5, Equals, clogger)

	clogger6 := blogger.GetChild("d.b.c")
	c.Assert(clogger6, Not(Equals), clogger)
}

func (s *LogriSuite) TestInheritLevelFromParent(chk *C) {
	a := s.logger.GetChild("a")
	b := s.logger.GetChild("a.b")
	c := s.logger.GetChild("a.b.c")
	d := s.logger.GetChild("a.b.c.d")
	e := s.logger.GetChild("a.b.c.d.e")

	s.logger.SetLevel(logrus.DebugLevel, true)
	b.SetLevel(logrus.ErrorLevel, false) // Don't propagate
	d.SetLevel(logrus.InfoLevel, true)

	s.AssertLogs(a.Info, "info")
	s.AssertLogs(a.Debug, "debug")

	s.AssertLogs(b.Error, "error")
	s.AssertNotLogs(b.Warn, "warn")

	s.AssertLogs(c.Debug, "debug")
	s.AssertLogs(c.Info, "info")

	s.AssertNotLogs(d.Debug, "debug")
	s.AssertLogs(d.Info, "info")

	s.AssertNotLogs(e.Debug, "debug")
	s.AssertLogs(e.Info, "info")

	// Unset d's level. Now d and e should inherit from the root, since b is a
	// non-propagate level
	d.SetLevel(NilLevel, true)
	s.AssertLogs(d.Debug, "debug")
	s.AssertLogs(e.Debug, "debug")

	// Set c's level to NilLevel, which it already is. Shouldn't affect anything
	c.SetLevel(NilLevel, true)
	s.AssertLogs(c.Debug, "debug")
	s.AssertLogs(d.Debug, "debug")

	// Now set c's level to something else and back to Nil. Still should
	// inherit from root
	c.SetLevel(logrus.FatalLevel, true)
	c.SetLevel(NilLevel, true)
	s.AssertLogs(c.Debug, "debug")
	s.AssertLogs(d.Debug, "debug")
	s.AssertLogs(e.Debug, "debug")
	s.AssertLogs(e.Debug, "debug")

}
