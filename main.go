package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework SystemConfiguration -framework CoreFoundation

// 作成したヘッダファイルをインクルードする
#include "monitor.h"
*/
import "C"

import (
	"fmt"
	"os"
	"runtime"
	"unsafe"
)

func main() {
	//OSスレッドに固定する
	runtime.LockOSThread()

	storeName := C.CString("GoNetworkMonitor")
	//終了時にメモリの解放
	defer C.free(unsafe.Pointer(storeName))

	//CF型にしてC言語で使えるようにエンコードする
	cfStoreName := C.CFStringCreateWithCString(
		C.CFAllocatorRef(unsafe.Pointer(nil)),
		storeName,
		C.kCFStringEncodingUTF8,
	)
	//ネットワークの監視の登録
	store := C.SCDynamicStoreCreate(
		C.CFAllocatorRef(unsafe.Pointer(nil)), //デフォルトメモリ管理
		cfStoreName,                           //名前
		(*[0]byte)(C.dynamicStoreCallback),    //状態が変化した時に呼び出す関数
		nil,
	)

	//storeがNull(登録できなかった）なら停止
	if store == (C.SCDynamicStoreRef)(unsafe.Pointer(nil)) {

		fmt.Println("Failed to create SCDynamicStore")
		os.Exit(1)
	}

	// モニタリング対象のキー（IPアドレス変更を検知）
	keys := []string{"State:/Network/Interface/.*/IPv4", "State:/Network/Interface/.*/IPv6"}
	cfKeys := makeCFArray(keys)
	//監視する内容の指定（cfKeysの中身の内容を監視する）
	C.SCDynamicStoreSetNotificationKeys(store, (C.CFArrayRef)(unsafe.Pointer(nil)), cfKeys)

	// RunLoopに登録
	rlSource := C.SCDynamicStoreCreateRunLoopSource(
		C.CFAllocatorRef(unsafe.Pointer(nil)),
		store,
		0,
	)

	//ホットラインの作成
	C.CFRunLoopAddSource(
		C.CFRunLoopGetCurrent(),
		rlSource,
		C.kCFRunLoopCommonModes,
	)

	fmt.Println("IP監視中...")
	//監視スタート
	C.CFRunLoopRun()
}

// Goの[]stringをCFArrayRefに変換
func makeCFArray(array []string) C.CFArrayRef {
	//配列を作成
	values := make([]C.CFStringRef, len(array))
	for i, s := range array {
		//c言語の文字列に変換
		cs := C.CString(s)
		values[i] = C.CFStringCreateWithCString(
			C.CFAllocatorRef(unsafe.Pointer(nil)),
			cs,
			C.kCFStringEncodingUTF8,
		)
		//メモリの解放
		C.free(unsafe.Pointer(cs))
	}
	return C.CFArrayCreate(
		C.CFAllocatorRef(unsafe.Pointer(nil)),
		(*unsafe.Pointer)(unsafe.Pointer(&values[0])),
		C.CFIndex(len(array)),
		nil,
	)
}
