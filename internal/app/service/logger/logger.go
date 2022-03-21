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

func (s *Service) Info(args ...interface{}) {
	s.stdOut.Info(args)
}

func (s *Service) Error(args ...interface{}) {
	s.stdErr.Error(args)
}
