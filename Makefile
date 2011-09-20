include $(GOROOT)/src/Make.inc

TARG=adp
GOFILES=\
	src/adp.go\
	src/console.go\
	src/labelDataSet.go\
	src/trainAndTest.go\
	src/convert.go\

include $(GOROOT)/src/Make.cmd
