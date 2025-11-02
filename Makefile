# Makefile for GVT
# Usage: make PLATFORM=linux|win
# Default PLATFORM is linux

PLATFORM ?= linux

ifeq ($(PLATFORM), linux)
include make/Makefile.linux
else ifeq ($(PLATFORM), win)
include make/Makefile.win
else
$(error "Unsupported PLATFORM '$(PLATFORM)'.")
endif