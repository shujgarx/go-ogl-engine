//go:build windows

package engine

/*
#cgo CFLAGS: -Wno-unused-variable
__declspec(dllexport) unsigned long NvOptimusEnablement = 0x00000001;
__declspec(dllexport) int AmdPowerXpressRequestHighPerformance = 1;
*/
import "C"
