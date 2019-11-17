package lib

import (
	"log"
)

type Logger struct {
	Err *log.Logger
	Inf *log.Logger
}

func (l *Logger) Error(err error){

	/*trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	l.Err.Output(2, trace)*/
	l.Inf.Println(err.Error())

}

func (l *Logger) Info(message string){

	l.Inf.Println(message)

}