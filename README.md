# GCP (go clang parse)
考虑到已有代码是 go，保持一致也写 go。

## BUILDING
要求：本地有 clang 10 已经编译。
编译了 clang, llvm-config, libclang。

```sh
$ ./llvm-config --libdir
/data/Programs/bytedance/llvmsrcs/build/lib
$ CGO_LDFLAGS="-L/data/Programs/bytedance/llvmsrcs/build/lib" go get -u github.com/go-clang/clang-v10/...
$ CGO_LDFLAGS="-L/data/Programs/bytedance/llvmsrcs/build/lib" go install github.com/go-clang/clang-v10/...
$ export LD_LIBRARY_PATH="/data/Programs/bytedance/llvmsrcs/build/lib"
$ export CGO_LDFLAGS="-L /data/Programs/bytedance/llvmsrcs/build/lib"
$ go build
$ ./gcp xxx.c
```

## LOGIC
main.go 读入 C 并解析。

关键逻辑
```
main:
	解析每个 C 文件得到一系列 translation unit `tu`
	对于每个 tu，在 transpile context 中分析 tu
	最后从 transpile context 总结得到 Repo 并 json 输出

TranspileContext.AddFile -> FileContext.analyze_tu

FileContext.analyze_tu
	利用 visitor pattern，用一个 cursor 访问 C 源码。
	分派任务给 visit_function
```

## RESTRICTIONS
1. 由于 C/CXX 无模块设计，现在为了简单所以只有一个 Package、只有一个 Module，都命名为 "."。也只有一个 repo。
2. parser 不处理 #define #include #if，它们会先被预处理器展开

## STATUS
* 文件数：单个
* 模块/包数目：单个
* 类型：只有 int
