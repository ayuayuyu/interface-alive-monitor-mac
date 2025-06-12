#ifndef MONITOR_H
#define MONITOR_H

// CGoで必要なフレームワークのヘッダをここにまとめる
#include <CoreFoundation/CoreFoundation.h>
#include <SystemConfiguration/SystemConfiguration.h>
#include <stdlib.h>

// Goから呼び出されるコールバック関数の宣言（プロトタイプ宣言）
void dynamicStoreCallback(SCDynamicStoreRef store, CFArrayRef changedKeys, void *info);

#endif // MONITOR_H