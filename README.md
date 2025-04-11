# DedupGo

DedupGo 是一个高效、安全的文件去重工具，提供直观的图形界面，帮助用户轻松管理和删除重复文件。

## ✨ 特性

- 🖥️ 现代化的图形界面
- 🔍 高效的文件扫描
- 🔐 支持多种哈希算法（MD5/SHA256）
- 📊 直观的重复文件展示
- 🗑️ 安全的文件删除（移动到回收站）
- 💻 跨平台支持（Windows/macOS/Linux）

## 🚀 安装

### 从源码编译

1. 确保已安装 Go 1.16 或更高版本
2. 克隆仓库：
```bash
git clone https://github.com/xiaozhe/dedupgo.git
cd dedupgo
```

3. 安装依赖：
```bash
go mod download
```

4. 编译项目：
```bash
# 编译命令行版本
go build -o dedupgo cmd/dedupgo/main.go

# 编译图形界面版本
go build -o dedupgo-gui cmd/dedupgo-gui/main.go
```

### 直接下载

访问 [Releases](https://github.com/xiaozhe/dedupgo/releases) 页面下载适合您系统的预编译版本。

## 📖 使用说明

### 图形界面版本

1. 启动程序
2. 点击"添加目录"选择要扫描的文件夹
3. 选择哈希算法（默认为 MD5）
4. 可选：设置最小文件大小过滤
5. 点击"开始扫描"
6. 查看扫描结果，选择性删除重复文件

### 命令行版本

```bash
# 基本用法
dedupgo scan /path/to/directory

# 使用 SHA256 算法
dedupgo scan -a sha256 /path/to/directory

# 设置最小文件大小（如：1MB）
dedupgo scan -s 1MB /path/to/directory
```

## 🛠️ 配置说明

### 支持的哈希算法
- MD5（默认）：速度快，适合一般使用
- SHA256：更高的安全性，但扫描速度较慢

### 文件大小过滤
- 支持的单位：KB、MB、GB
- 示例：1MB、500KB、2GB

## 🔒 安全性说明

- 重复文件删除时会移动到回收站而不是直接删除
- 始终保留一个原始文件，不会删除所有副本
- 使用可靠的哈希算法确保文件比对准确性

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

本项目采用 MIT 许可证，详见 [LICENSE](LICENSE) 文件。

## 🙏 致谢

- [Fyne](https://fyne.io/) - 跨平台 GUI 框架
- 所有贡献者和用户

## 📞 联系方式

如有问题或建议，欢迎通过以下方式联系：

- 提交 [Issue](https://github.com/xiaozhe/dedupgo/issues)
- 邮件：xiaozhe9629@gmail.com
