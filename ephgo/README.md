# Swiss Ephemeris Go版本 (EphGo)

这是一个将瑞士星历表（Swiss Ephemeris）从C语言转换为Go语言的实现，专注于支持JPL文件格式读取和精确计算行星、月球、恒星等天体的位置。

## 特性

- 🌟 **精确的天体位置计算** - 支持太阳系所有主要天体
- 📅 **完整的日期时间处理** - 儒略历与格里高利历互转
- 🛸 **JPL星历文件支持** - 读取和解析JPL DE系列星历文件
- 🌙 **月球交点和远地点** - 计算月球轨道特殊点
- 📊 **多种坐标系统** - 支持地心、日心、极坐标、笛卡尔坐标
- 🔧 **线程安全** - 支持并发计算
- 🚀 **高性能** - Go语言原生实现，避免CGO开销

## 安装

```bash
go get ephgo
```

## 快速开始

```go
package main

import (
    "fmt"
    "ephgo"
)

func main() {
    // 设置星历文件路径（可选）
    ephgo.SetEphePath("./ephe")
    
    // 计算当前时间太阳的位置
    tjd := ephgo.GetCurrentTime()
    xx, err := ephgo.CalcUT(tjd, ephgo.SeSun, ephgo.SeflgSwieph)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("太阳位置: 经度=%.6f°, 纬度=%.6f°, 距离=%.6f AU\n", 
        xx[0], xx[1], xx[2])
    
    // 清理资源
    ephgo.Close()
}
```

## 主要模块

### 1. 常量定义 (constants.go)
- 天体编号和名称
- 物理常数
- 计算标志位
- 数学常数

### 2. 数据类型 (types.go)
- 星历文件数据结构
- 行星数据结构
- 全局状态管理
- 线程安全访问

### 3. 日期时间 (date.go)
- 儒略日计算
- 历法转换
- UTC/ET时间转换
- Delta T计算

### 4. JPL星历 (jpl.go)
- JPL文件读取
- 切比雪夫多项式插值
- 天体位置计算
- 坐标系转换

### 5. 核心计算 (sweph.go)
- 主要计算接口
- 天体位置计算
- 坐标变换
- 星历数据管理

## 支持的天体

### 主要行星
- 太阳 (SeSun)
- 月球 (SeMoon)
- 水星 (SeMercury)
- 金星 (SeVenus)
- 火星 (SeMars)
- 木星 (SeJupiter)
- 土星 (SeSaturn)
- 天王星 (SeUranus)
- 海王星 (SeNeptune)
- 冥王星 (SePluto)
- 地球 (SeEarth)

### 月球轨道点
- 平交点 (SeMeanNode)
- 真交点 (SeTrueNode)
- 平远地点 (SeMeanApog)
- 振荡远地点 (SeOscuApog)

### 小行星
- 谷神星 (SeCeres)
- 智神星 (SePallas)
- 婚神星 (SeJuno)
- 灶神星 (SeVesta)

### 小天体
- 凯龙星 (SeChiron)
- 福鲁斯 (SePholus)

## 计算标志位

```go
// 星历类型
SeflgJpleph   // 使用JPL星历
SeflgSwieph   // 使用Swiss Ephemeris星历
SeflgMoseph   // 使用Moshier星历

// 坐标系统
SeflgHelctr   // 日心坐标
SeflgBaryctr  // 重心坐标
SeflgTopoctr  // 地心坐标

// 坐标格式
SeflgXyz      // 笛卡尔坐标
SeflgRadians  // 弧度单位
SeflgEquatorial // 赤道坐标

// 精度选项
SeflgSpeed    // 计算速度
SeflgTruepos  // 真实位置
SeflgNoaberr  // 无光行差
SeflgNogdefl  // 无引力偏折
```

## 使用示例

### 计算行星位置

```go
// 2024年1月1日12:00 UTC
tjd := ephgo.Julday(2024, 1, 1, 12.0, ephgo.SeGregCal)

// 计算火星位置（地心坐标）
xx, err := ephgo.Calc(tjd, ephgo.SeMars, ephgo.SeflgSwieph)
if err != nil {
    panic(err)
}

fmt.Printf("火星位置:\n")
fmt.Printf("  经度: %.6f°\n", xx[0])
fmt.Printf("  纬度: %.6f°\n", xx[1])
fmt.Printf("  距离: %.6f AU\n", xx[2])
```

### 不同坐标系统

```go
tjd := ephgo.Julday(2024, 1, 1, 12.0, ephgo.SeGregCal)

// 地心坐标（度）
xx1, _ := ephgo.Calc(tjd, ephgo.SeMoon, ephgo.SeflgSwieph)

// 地心笛卡尔坐标
xx2, _ := ephgo.Calc(tjd, ephgo.SeMoon, ephgo.SeflgSwieph|ephgo.SeflgXyz)

// 日心坐标
xx3, _ := ephgo.Calc(tjd, ephgo.SeMoon, ephgo.SeflgSwieph|ephgo.SeflgHelctr)

// 弧度单位
xx4, _ := ephgo.Calc(tjd, ephgo.SeMoon, ephgo.SeflgSwieph|ephgo.SeflgRadians)
```

### 日期时间转换

```go
// 格里高利历到儒略日
jd := ephgo.Julday(2024, 1, 1, 12.0, ephgo.SeGregCal)

// 儒略日到格里高利历
year, month, day, jut := ephgo.Revjul(jd, ephgo.SeGregCal)

// UTC到儒略日（包含Delta T转换）
et, ut1, err := ephgo.UtcToJd(2024, 1, 1, 12, 0, 0.0, ephgo.SeGregCal)

// 计算星期几
dayOfWeek := ephgo.DayOfWeek(jd)

// 判断闰年
isLeap := ephgo.IsLeapYear(2024)
```

### 月球交点和远地点

```go
tjd := ephgo.Julday(2024, 1, 1, 12.0, ephgo.SeGregCal)

// 月球平交点
meanNode, _ := ephgo.Calc(tjd, ephgo.SeMeanNode, ephgo.SeflgSwieph)

// 月球真交点
trueNode, _ := ephgo.Calc(tjd, ephgo.SeTrueNode, ephgo.SeflgSwieph)

// 月球平远地点
meanApogee, _ := ephgo.Calc(tjd, ephgo.SeMeanApog, ephgo.SeflgSwieph)
```

## 文件结构

```
ephgo/
├── go.mod              # Go模块定义
├── constants.go        # 常量定义
├── types.go           # 数据类型
├── date.go            # 日期时间处理
├── jpl.go             # JPL文件读取
├── sweph.go           # 核心计算
├── example.go         # 示例程序
└── README.md          # 项目说明
```

## 精度和性能

- **位置精度**: 角秒级别（取决于使用的星历文件）
- **时间范围**: 公元前13200年到公元17191年（DE431）
- **计算速度**: 单次计算通常在微秒级别
- **内存使用**: 最小化内存占用，支持大量并发计算

## 星历文件

该库支持JPL DE系列星历文件：
- DE405: 1600-2200年，高精度
- DE406: 3000 BC - 3000 AD，长期
- DE430: 1550-2650年，最新月球模型
- DE431: -13200 - +17191年，超长期

星历文件可从以下来源获取：
- [JPL官方网站](https://ssd.jpl.nasa.gov/)
- [Swiss Ephemeris网站](https://www.astro.com/swisseph/)

## 注意事项

1. **星历文件路径**: 确保设置正确的星历文件路径
2. **时间范围**: 检查计算时间是否在星历文件覆盖范围内
3. **坐标系统**: 明确指定所需的坐标系统和单位
4. **资源清理**: 程序结束时调用 `ephgo.Close()` 释放资源
5. **错误处理**: 始终检查函数返回的错误信息

## 限制

当前版本的限制：
- 主要支持JPL星历文件，Swiss Ephemeris二进制文件支持有限
- Moshier理论星历尚未完全实现
- 一些高级功能（如恒星、小行星详细计算）仍在开发中

## 示例程序

运行示例程序：

```bash
cd ephgo
go run example.go
```

示例程序演示了：
- 当前时间太阳位置计算
- 指定日期行星位置计算
- 不同坐标系统使用
- 月球交点和远地点计算
- 日期时间转换功能

## 贡献

欢迎提交问题报告和功能请求。如果要贡献代码，请：

1. Fork 项目
2. 创建功能分支
3. 提交更改
4. 创建 Pull Request

## 许可证

本项目基于 AGPL-3.0 许可证开源，与原始Swiss Ephemeris保持一致。

## 致谢

- Swiss Ephemeris项目的原作者 Dieter Koch 和 Alois Treindl
- JPL开发计算中心提供的精确星历数据
- Go语言社区的支持和工具

## 联系方式

如有问题或建议，请通过以下方式联系：
- 提交 GitHub Issue
- 发送邮件至项目维护者

---

**免责声明**: 本软件仅供学习和研究使用。虽然我们努力确保计算精度，但不对任何特定用途的适用性或精度做出保证。在关键应用中使用前，请验证计算结果。