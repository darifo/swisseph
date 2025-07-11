package main

import (
	"fmt"
	"log"
	"time"

	"ephgo"
)

func main() {
	fmt.Println("Swiss Ephemeris Go版本示例程序")
	fmt.Printf("版本: %s\n\n", ephgo.Version())

	// 设置星历文件路径（可选）
	ephgo.SetEphePath("./ephe")

	// 示例1：计算当前时间太阳的位置
	fmt.Println("=== 示例1：当前时间太阳位置 ===")
	calculateCurrentSunPosition()

	// 示例2：计算指定时间的行星位置
	fmt.Println("\n=== 示例2：2024年1月1日行星位置 ===")
	calculatePlanetsOn20240101()

	// 示例3：使用不同的坐标系
	fmt.Println("\n=== 示例3：不同坐标系示例 ===")
	demonstrateCoordinateSystems()

	// 示例4：月球交点和远地点
	fmt.Println("\n=== 示例4：月球交点和远地点 ===")
	calculateLunarNodesAndApogee()

	// 示例5：日期时间转换示例
	fmt.Println("\n=== 示例5：日期时间转换 ===")
	demonstrateDateTimeConversion()

	// 清理资源
	ephgo.Close()
	fmt.Println("\n程序结束")
}

// calculateCurrentSunPosition 计算当前时间太阳的位置
func calculateCurrentSunPosition() {
	// 获取当前时间的儒略日
	tjdUT := ephgo.GetCurrentTime()
	
	// 计算太阳位置（地心坐标，视位置）
	xx, err := ephgo.CalcUT(tjdUT, ephgo.SeSun, ephgo.SeflgSwieph)
	if err != nil {
		log.Printf("计算太阳位置失败: %v", err)
		return
	}
	
	fmt.Printf("当前UTC时间: %s\n", time.Now().UTC().Format("2006-01-02 15:04:05"))
	fmt.Printf("儒略日 (UT): %.6f\n", tjdUT)
	fmt.Printf("太阳位置:\n")
	fmt.Printf("  经度: %.6f°\n", xx[0])
	fmt.Printf("  纬度: %.6f°\n", xx[1])
	fmt.Printf("  距离: %.6f AU\n", xx[2])
	fmt.Printf("  经度速度: %.6f°/日\n", xx[3])
	fmt.Printf("  纬度速度: %.6f°/日\n", xx[4])
	fmt.Printf("  距离速度: %.6f AU/日\n", xx[5])
}

// calculatePlanetsOn20240101 计算2024年1月1日的行星位置
func calculatePlanetsOn20240101() {
	// 2024年1月1日 12:00 UTC的儒略日
	tjd := ephgo.Julday(2024, 1, 1, 12.0, ephgo.SeGregCal)
	
	// 行星列表
	planets := []struct {
		id   int
		name string
	}{
		{ephgo.SeSun, "太阳"},
		{ephgo.SeMoon, "月球"},
		{ephgo.SeMercury, "水星"},
		{ephgo.SeVenus, "金星"},
		{ephgo.SeMars, "火星"},
		{ephgo.SeJupiter, "木星"},
		{ephgo.SeSaturn, "土星"},
		{ephgo.SeUranus, "天王星"},
		{ephgo.SeNeptune, "海王星"},
		{ephgo.SePluto, "冥王星"},
	}
	
	fmt.Printf("日期: 2024年1月1日 12:00 UTC (JD %.6f)\n", tjd)
	fmt.Println("天体          经度(°)     纬度(°)     距离(AU)")
	fmt.Println("----------------------------------------")
	
	for _, planet := range planets {
		xx, err := ephgo.Calc(tjd, planet.id, ephgo.SeflgSwieph)
		if err != nil {
			fmt.Printf("%-8s: 计算失败 - %v\n", planet.name, err)
			continue
		}
		
		fmt.Printf("%-8s  %9.4f  %9.4f  %9.6f\n", 
			planet.name, xx[0], xx[1], xx[2])
	}
}

// demonstrateCoordinateSystems 演示不同坐标系
func demonstrateCoordinateSystems() {
	tjd := ephgo.Julday(2024, 1, 1, 12.0, ephgo.SeGregCal)
	
	fmt.Printf("日期: 2024年1月1日 12:00 UTC\n")
	fmt.Printf("天体: 月球\n\n")
	
	// 地心坐标（默认）
	xx1, err := ephgo.Calc(tjd, ephgo.SeMoon, ephgo.SeflgSwieph)
	if err == nil {
		fmt.Printf("地心坐标 (度):\n")
		fmt.Printf("  经度: %.6f°, 纬度: %.6f°, 距离: %.6f AU\n", xx1[0], xx1[1], xx1[2])
	}
	
	// 地心笛卡尔坐标
	xx2, err := ephgo.Calc(tjd, ephgo.SeMoon, ephgo.SeflgSwieph|ephgo.SeflgXyz)
	if err == nil {
		fmt.Printf("地心笛卡尔坐标 (AU):\n")
		fmt.Printf("  X: %.6f, Y: %.6f, Z: %.6f\n", xx2[0], xx2[1], xx2[2])
	}
	
	// 弧度单位
	xx3, err := ephgo.Calc(tjd, ephgo.SeMoon, ephgo.SeflgSwieph|ephgo.SeflgRadians)
	if err == nil {
		fmt.Printf("地心坐标 (弧度):\n")
		fmt.Printf("  经度: %.6f rad, 纬度: %.6f rad, 距离: %.6f AU\n", xx3[0], xx3[1], xx3[2])
	}
	
	// 日心坐标
	xx4, err := ephgo.Calc(tjd, ephgo.SeMoon, ephgo.SeflgSwieph|ephgo.SeflgHelctr)
	if err == nil {
		fmt.Printf("日心坐标 (度):\n")
		fmt.Printf("  经度: %.6f°, 纬度: %.6f°, 距离: %.6f AU\n", xx4[0], xx4[1], xx4[2])
	}
}

// calculateLunarNodesAndApogee 计算月球交点和远地点
func calculateLunarNodesAndApogee() {
	tjd := ephgo.Julday(2024, 1, 1, 12.0, ephgo.SeGregCal)
	
	fmt.Printf("日期: 2024年1月1日 12:00 UTC\n\n")
	
	// 月球平交点
	xxMeanNode, err := ephgo.Calc(tjd, ephgo.SeMeanNode, ephgo.SeflgSwieph)
	if err == nil {
		fmt.Printf("月球平交点:\n")
		fmt.Printf("  经度: %.6f°\n", xxMeanNode[0])
		fmt.Printf("  速度: %.6f°/日\n", xxMeanNode[3])
	} else {
		fmt.Printf("月球平交点计算失败: %v\n", err)
	}
	
	// 月球真交点
	xxTrueNode, err := ephgo.Calc(tjd, ephgo.SeTrueNode, ephgo.SeflgSwieph)
	if err == nil {
		fmt.Printf("\n月球真交点:\n")
		fmt.Printf("  经度: %.6f°\n", xxTrueNode[0])
		fmt.Printf("  速度: %.6f°/日\n", xxTrueNode[3])
	} else {
		fmt.Printf("\n月球真交点计算失败: %v\n", err)
	}
	
	// 月球平远地点
	xxMeanApog, err := ephgo.Calc(tjd, ephgo.SeMeanApog, ephgo.SeflgSwieph)
	if err == nil {
		fmt.Printf("\n月球平远地点:\n")
		fmt.Printf("  经度: %.6f°\n", xxMeanApog[0])
		fmt.Printf("  速度: %.6f°/日\n", xxMeanApog[3])
	} else {
		fmt.Printf("\n月球平远地点计算失败: %v\n", err)
	}
}

// demonstrateDateTimeConversion 演示日期时间转换
func demonstrateDateTimeConversion() {
	// 格里高利历到儒略日转换
	jd := ephgo.Julday(2024, 1, 1, 12.5, ephgo.SeGregCal)
	fmt.Printf("2024年1月1日 12:30 UTC = JD %.6f\n", jd)
	
	// 儒略日到格里高利历转换
	year, month, day, jut := ephgo.Revjul(jd, ephgo.SeGregCal)
	hours := int(jut)
	minutes := int((jut - float64(hours)) * 60)
	seconds := ((jut - float64(hours)) * 60 - float64(minutes)) * 60
	
	fmt.Printf("JD %.6f = %d年%d月%d日 %02d:%02d:%06.3f UTC\n", 
		jd, year, month, day, hours, minutes, seconds)
	
	// UTC到儒略日转换
	et, ut1, err := ephgo.UtcToJd(2024, 1, 1, 12, 30, 0.0, ephgo.SeGregCal)
	if err == nil {
		fmt.Printf("2024-01-01 12:30:00 UTC:\n")
		fmt.Printf("  ET = JD %.6f\n", et)
		fmt.Printf("  UT1 = JD %.6f\n", ut1)
		fmt.Printf("  Delta T = %.3f 秒\n", (et-ut1)*86400)
	}
	
	// 星期几计算
	dayOfWeek := ephgo.DayOfWeek(jd)
	weekdays := []string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}
	fmt.Printf("2024年1月1日是%s\n", weekdays[dayOfWeek])
	
	// 闰年判断
	isLeap := ephgo.IsLeapYear(2024)
	fmt.Printf("2024年是闰年: %t\n", isLeap)
	
	// 某月天数
	daysInFeb := ephgo.DaysInMonth(2024, 2)
	fmt.Printf("2024年2月有%d天\n", daysInFeb)
}

// formatDegrees 格式化角度显示
func formatDegrees(deg float64) string {
	d := int(deg)
	m := int((deg - float64(d)) * 60)
	s := ((deg - float64(d)) * 60 - float64(m)) * 60
	return fmt.Sprintf("%d°%02d'%05.2f\"", d, m, s)
}

// formatTime 格式化时间显示
func formatTime(jd float64) string {
	year, month, day, jut := ephgo.Revjul(jd, ephgo.SeGregCal)
	hours := int(jut)
	minutes := int((jut - float64(hours)) * 60)
	seconds := ((jut - float64(hours)) * 60 - float64(minutes)) * 60
	
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%06.3f", 
		year, month, day, hours, minutes, seconds)
}