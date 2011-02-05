include $(GOROOT)/src/Make.inc

TARG=adp
GOFILES=\
	adp.go\
	labelDataSet.go\
	trainAndTest.go\

include $(GOROOT)/src/Make.cmd
