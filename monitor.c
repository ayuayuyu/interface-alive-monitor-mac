#include "monitor.h" // 作成したヘッダファイルを読み込む
#include <stdio.h>   // printfを使うために追加

// コールバック関数の具体的な処理内容（実装）
void dynamicStoreCallback(SCDynamicStoreRef store, CFArrayRef changedKeys, void *info)
{
    CFIndex count = CFArrayGetCount(changedKeys);
    for (CFIndex i = 0; i < count; i++)
    {
        CFStringRef key = (CFStringRef)CFArrayGetValueAtIndex(changedKeys, i);
        char buf[1024];
        if (CFStringGetCString(key, buf, sizeof(buf), kCFStringEncodingUTF8))
        {
            // Cのコードなので、Goのfmt.PrintlnではなくCのprintfを使う
            printf("Changed key: %s\n", buf);
        }
    }
}