# 故障排除指南

## 编译问题：C编译器未找到

### 问题描述

在尝试运行 Go 程序时，遇到以下错误：

```
cgo: C compiler "/opt/homebrew/Cellar/gcc/14.1.0_1/bin/gcc-14" not found: exec: "/opt/homebrew/Cellar/gcc/14.1.0_1/bin/gcc-14": stat /opt/homebrew/Cellar/gcc/14.1.0_1/bin/gcc-14: no such file or directory
```

### 原因

这个错误是由于系统无法找到指定的 GCC 编译器造成的。在 macOS 系统上，默认的 C 编译器通常是 Clang，而不是 GCC。

### 解决方案

1. 首先，检查系统中已安装的 C 编译器：

   ```
   gcc --version
   ```

   在 macOS 上，这通常会显示 Clang 的版本信息。

2. 使用系统默认的 Clang 编译器来编译和运行 Go 程序：

   ```
   CGO_ENABLED=1 CC=clang go run cmd/api/main.go
   ```

   这个命令做了以下事情：
   - 设置 `CGO_ENABLED=1` 以启用 cgo
   - 设置 `CC=clang` 以使用系统默认的 Clang 编译器
   - 运行 Go 程序

### 注意事项

- 如果在未来的开发中需要使用特定版本的 GCC，可以考虑使用 Homebrew 安装并管理 GCC 版本。
- 在项目文档中记录这些环境设置，以确保团队成员都能正确编译和运行项目。
