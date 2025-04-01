# GCP (go clang parse)
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

