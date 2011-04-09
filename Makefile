include $(GOROOT)/src/Make.inc

TARG=adp
SRC_DIR=src
GOFILES=\
	$(SRC_DIR)/adp.go\
	$(SRC_DIR)/labelDataSet.go\
	$(SRC_DIR)/trainAndTest.go\

include $(GOROOT)/src/Make.cmd
