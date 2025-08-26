# getNoPSS - Kubernetes Pod Security Scanner

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![OpenAI](https://img.shields.io/badge/OpenAI-GPT4o-orange.svg)](https://openai.com)

🔒 **getNoPSS** 是一个强大的 Kubernetes Pod 安全扫描工具，结合传统安全检查和 AI 驱动的智能分析，帮助识别和修复容器安全问题。

## ✨ 特性

### 🛡️ 传统安全检查
- **Pod Security Standards (PSS) 合规性检查**
- **13 种安全风险检测**：特权容器、主机网络、危险挂载、权限提升等
- **命名空间过滤**：灵活的扫描范围控制
- **详细报告**：清晰的安全问题和位置信息

### 🤖 AI 智能分析
- **深度安全分析**：使用 OpenAI GPT-4o 进行上下文感知的安全评估
- **智能风险评级**：SAFE、MODERATE、HIGH_RISK、CRITICAL 四级风险分类
- **个性化建议**：针对每个 Pod 的具体安全改进建议
- **多种输出格式**：JSON、HTML、控制台输出

### 🔧 易用性
- **简单配置**：YAML 配置文件，无需复杂设置
- **API 兼容性**：支持官方 OpenAI 和兼容的第三方 API
- **详细文档**：完整的使用指南和示例

## 📦 安装

### 从源码构建

```bash
# 克隆项目
git clone https://github.com/yourusername/getNoPSS.git
cd getNoPSS

# 构建
go build -o getNoPSS main.go
```

### 系统要求

- Go 1.21+
- 有效的 kubeconfig 配置
- OpenAI API 密钥（用于 AI 分析功能）

## 🚀 快速开始

### 1. 生成配置文件

```bash
./getNoPSS generateConfig
```

### 2. 编辑配置文件

编辑生成的 `config.yaml` 文件：

```yaml
openai:
  api_key: "your-openai-api-key"
  base_url: "https://api.openai.com/v1"  # 或使用兼容的第三方API
  model: "gpt-4o"
```

### 3. 测试 API 连接

```bash
./getNoPSS testApi
```

### 4. 运行安全扫描

```bash
# 传统安全检查
./getNoPSS allNoPSS

# AI 智能分析
./getNoPSS aiAnalysis -c
```

## 📋 使用指南

### 传统安全检查

```bash
# 扫描所有 Pod
./getNoPSS allNoPSS

# 排除系统命名空间
./getNoPSS allNoPSS -e kube-system,kube-public
```

### AI 智能分析

```bash
# 基础 AI 分析
./getNoPSS aiAnalysis

# 生成 HTML 报告
./getNoPSS aiAnalysis -f html -o security_report.html

# 控制台显示详细结果
./getNoPSS aiAnalysis -c

# 使用自定义配置文件
./getNoPSS aiAnalysis --config my-config.yaml

# 排除特定命名空间
./getNoPSS aiAnalysis -e kube-system,kube-public -c
```

## 🔍 检测的安全问题

### 传统检查项目

1. **Host PID** - 使用主机 PID 命名空间
2. **Host Network** - 使用主机网络
3. **Host IPC** - 使用主机 IPC 命名空间
4. **Host Ports** - 使用主机端口
5. **Host Path** - 挂载主机路径卷
6. **Host Process** - Windows HostProcess 容器
7. **Privileged** - 特权容器
8. **Allow Privilege Escalation** - 允许权限提升
9. **Added Capabilities** - 添加的 Linux 权限
10. **Dropped Capabilities** - 移除的 Linux 权限
11. **Seccomp Disabled** - 禁用 Seccomp
12. **AppArmor Disabled** - 禁用 AppArmor
13. **Unmasked Procmount** - 未屏蔽的 proc 挂载
14. **Unsafe Sysctl** - 不安全的 sysctl 设置

### AI 分析维度

- 🔐 **安全上下文配置**
- 📦 **镜像安全性**
- 🌐 **网络策略**
- 💾 **资源限制**
- 🔑 **权限管理**
- 📊 **合规性检查**

## 📊 输出格式

### JSON 格式
```json
{
  "generated_at": "2024-08-22T08:00:00Z",
  "total_pods": 10,
  "summary": {
    "SAFE": 3,
    "MODERATE": 4,
    "HIGH_RISK": 2,
    "CRITICAL": 1
  },
  "analyses": [...]
}
```
<img width="960" height="721" alt="image" src="https://github.com/user-attachments/assets/6f939f12-2381-4b85-86ba-0d860acd5672" />

### HTML 报告
- 📈 统计概览
- 🎨 颜色编码的风险级别
- 📝 详细的问题描述和建议
- 📱 响应式设计

## ⚙️ 配置选项

### 配置文件结构

```yaml
openai:
  # OpenAI API 密钥 (必需)
  api_key: "your-api-key"
  
  # API Base URL (可选)
  # 官方: https://api.openai.com/v1
  # 代理: https://api.example.com/v1
  base_url: ""
  
  # 模型选择 (可选，默认: gpt-4o)
  model: "gpt-4o"
```

### 命令行参数

| 参数 | 描述 | 默认值 |
|------|------|--------|
| `--config` | 配置文件路径 | `config.yaml` |
| `-o, --output` | 输出文件路径 | 自动生成 |
| `-f, --format` | 输出格式 (json\|html) | `json` |
| `-c, --console` | 控制台显示详细结果 | `false` |
| `-e, --exclude` | 排除的命名空间列表 | - |

## 🤝 贡献

欢迎贡献代码！请遵循以下步骤：

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 📝 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🆘 支持

- 📖 [文档](CLAUDE.md)
- 🐛 [问题反馈](https://github.com/yourusername/getNoPSS/issues)
- 💬 [讨论](https://github.com/yourusername/getNoPSS/discussions)

## 🙏 致谢

- [Kubernetes](https://kubernetes.io/) - 容器编排平台
- [OpenAI](https://openai.com/) - AI 能力支持
- [Cobra](https://github.com/spf13/cobra) - CLI 框架
- [Zerolog](https://github.com/rs/zerolog) - 高性能日志库

---

⭐ 如果这个项目对你有帮助，请给个 Star！
