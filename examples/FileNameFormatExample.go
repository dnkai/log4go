package main

import (
	"github.com/log4go"
	"time"
)

func main()  {
	log4go.LoadConfiguration("./example.xml")
	//flw := log4go.NewFileLogWriter("tt-%P.log", true)
	//flw.SetFormat("%A %B [%D %T] [%L] (%S) %M")
	//flw.SetModuleName("BTC")
	//flw.SetOutputMode(log4go.OutputModeJson)

	flw1 := log4go.NewFileLogWriter("log-%P.log", true)
	flw1.SetModuleName("ETH")

	flw2 := log4go.NewJSONLogWriter("log-%P.log", true)
	flw2.SetModuleName("XLM")

	l1 := log4go.NewDefaultLogger(log4go.DEBUG)
	l2 := log4go.NewDefaultLogger(log4go.DEBUG)

	defer func() {
		log4go.Close()
		l1.Close()
		l2.Close()
		time.Sleep(time.Second)
	}()

	//log4go.AddFilter("tt", log4go.DEBUG, flw)
	l1.AddFilter("tt", log4go.DEBUG, flw1)
	l2.AddFilter("tt", log4go.DEBUG, flw2)

	// 1
	log4go.Trace("this is trace")
	log4go.Info("this is info")
	log4go.Error("this is error")
	l1.Trace("this is trace")
	l1.Info("this is info")
	l1.Error("this is error")
	l2.Trace("this is trace")
	l2.Info("this is info")
	l2.Error("this is error")

	//
	//flw.Rotate()
	flw1.Rotate()
	flw2.Rotate()
	log4go.Trace("this is trace-2")
	log4go.Info("this is info-2")
	log4go.Error("this is error-2")
	l1.Trace("this is trace-2")
	l1.Info("this is info-2")
	l1.Error("this is error-2")
	l2.Trace("this is trace-2")
	l2.Info("this is info-2")
	l2.Error("this is error-2")
}