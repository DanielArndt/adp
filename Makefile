include $(GOROOT)/src/Make.inc

TARG=adp
<<<<<<< HEAD
SRC_DIR=src
GOFILES=\
	$(SRC_DIR)/adp.go\
	$(SRC_DIR)/misc.go\
	$(SRC_DIR)/labelDataSet.go\
	$(SRC_DIR)/trainAndTest.go\
=======
GOFILES=\
	src/adp.go\
	src/console.go\
	src/labelDataSet.go\
	src/trainAndTest.go\
	src/convert.go\
>>>>>>> multiTestSet

include $(GOROOT)/src/Make.cmd
