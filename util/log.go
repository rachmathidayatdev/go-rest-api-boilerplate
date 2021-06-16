package util

import "github.com/sirupsen/logrus"

//Logger type
type Logger logrus.FieldLogger

//Log var
var Log logrus.FieldLogger

//BuildContext func
func BuildContext(context string) logrus.Fields {
	return logrus.Fields{
		"context": context,
	}
}
