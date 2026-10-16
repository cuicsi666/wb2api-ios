//go:build carchive

// 通过 -buildmode=c-archive 把 workbuddy2api 编译成 iOS 可链接的静态库时使用的 C 入口。
//
// 构建（macOS + Xcode）：
//
//	GOOS=ios GOARCH=arm64 CGO_ENABLED=1 \
//	CC="$(xcrun --sdk iphoneos -f clang)" \
//	CGO_CFLAGS="-isysroot $(xcrun --sdk iphoneos --show-sdk-path) -miphoneos-version-min=15.0 -arch arm64" \
//	CGO_LDFLAGS="-isysroot $(xcrun --sdk iphoneos --show-sdk-path) -miphoneos-version-min=15.0 -arch arm64" \
//	go build -tags carchive -buildmode=c-archive -o libwb.a .
package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"

	"workbuddy2api/wbios"
)

// WbStart 启动网关；成功返回 NULL，失败返回错误信息（须用 WbFree 释放）。
//
//export WbStart
func WbStart(cDataDir, cListen, cAPIKey *C.char) *C.char {
	if err := wbios.Start(C.GoString(cDataDir), C.GoString(cListen), C.GoString(cAPIKey)); err != nil {
		return C.CString(err.Error())
	}
	return nil
}

// WbStop 停止网关。
//
//export WbStop
func WbStop() *C.char {
	if err := wbios.Stop(); err != nil {
		return C.CString(err.Error())
	}
	return nil
}

// WbStatus 返回状态快照 JSON（须用 WbFree 释放）。
//
//export WbStatus
func WbStatus() *C.char { return C.CString(wbios.StatusJSON()) }

// WbVersion 版本号（须用 WbFree 释放）。
//
//export WbVersion
func WbVersion() *C.char { return C.CString(wbios.Version) }

// WbAddr 监听地址（须用 WbFree 释放）。
//
//export WbAddr
func WbAddr() *C.char { return C.CString(wbios.Addr()) }

// WbAuthDir 凭证目录（须用 WbFree 释放）。
//
//export WbAuthDir
func WbAuthDir() *C.char { return C.CString(wbios.AuthDir()) }

// WbRunning 1=运行中，0=已停止。
//
//export WbRunning
func WbRunning() C.int {
	if wbios.Running() {
		return 1
	}
	return 0
}

// WbUptime 已运行秒数。
//
//export WbUptime
func WbUptime() C.int { return C.int(wbios.Uptime()) }

// WbFree 释放由本库返回的 C 字符串。
//
//export WbFree
func WbFree(p *C.char) {
	if p != nil {
		C.free(unsafe.Pointer(p))
	}
}

// main 不会被调用；c-archive 要求包为 main。
func main() {}
