package logger

import (
	"github.com/sirupsen/logrus"
)

type Service struct {
	stdOut *logrus.Logger
	stdErr *logrus.Logger
}

func NewLoggerService(stdOut *logrus.Logger, stdErr *logrus.Logger) *Service {

	return &Service{
		stdOut: stdOut,
		stdErr: stdErr,
	}
}

//nolint
func (s *Service) Info(args ...interface{}) {
	s.stdOut.Info(args)
}

//nolint
func (s *Service) Error(args ...interface{}) {
	s.stdErr.Error(args)
}
