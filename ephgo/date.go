package ephgo

import (
	"fmt"
	"math"
	"time"
)

// DateConversion 日期转换函数
// y: 年, m: 月, d: 日, utime: 世界时（小时，十进制）
// c: 历法类型 'g'(格里高利历) 或 'j'(儒略历)
// 返回格林威治标准时（天）
func DateConversion(y, m, d int, utime Float64, c byte) (Float64, error) {
	var gregflag int
	if c == 'g' || c == 'G' {
		gregflag = SeGregCal
	} else if c == 'j' || c == 'J' {
		gregflag = SeJulCal
	} else {
		return 0, fmt.Errorf("无效的历法标志: %c", c)
	}
	
	jd := Julday(y, m, d, utime, gregflag)
	return jd, nil
}

// Julday 计算儒略日
// year: 年, month: 月, day: 日, hour: 小时（十进制）
// gregflag: 历法标志（SeGregCal或SeJulCal）
func Julday(year, month, day int, hour Float64, gregflag int) Float64 {
	var jd Float64
	var a, b int
	
	if month <= 2 {
		year--
		month += 12
	}
	
	if gregflag == SeGregCal {
		a = year / 100
		b = 2 - a + a/4
	} else {
		b = 0
	}
	
	if year < 0 {
		jd = math.Floor(365.25*Float64(year-1)) + math.Floor(30.6001*Float64(month+1)) + Float64(day) + Float64(b) + 1720994.5
	} else {
		jd = math.Floor(365.25*Float64(year)) + math.Floor(30.6001*Float64(month+1)) + Float64(day) + Float64(b) + 1720994.5
	}
	
	jd += hour / 24.0
	return jd
}

// Revjul 从儒略日反推历法日期
// jd: 儒略日
// gregflag: 历法标志
// 返回：年、月、日、世界时（小时）
func Revjul(jd Float64, gregflag int) (year, month, day int, jut Float64) {
	var a, b, c, d, e int
	var f Float64
	
	jd += 0.5
	z := int(jd)
	f = jd - Float64(z)
	
	if gregflag == SeGregCal {
		if z >= 2299161 {
			alpha := int((Float64(z) - 1867216.25) / 36524.25)
			a = z + 1 + alpha - alpha/4
		} else {
			a = z
		}
	} else {
		a = z
	}
	
	b = a + 1524
	c = int((Float64(b) - 122.1) / 365.25)
	d = int(365.25 * Float64(c))
	e = int(Float64(b-d) / 30.6001)
	
	day = b - d - int(30.6001*Float64(e))
	
	if e < 14 {
		month = e - 1
	} else {
		month = e - 13
	}
	
	if month > 2 {
		year = c - 4716
	} else {
		year = c - 4715
	}
	
	jut = f * 24.0
	return
}

// UtcToJd 将UTC时间转换为儒略日
// iyear, imonth, iday: 年月日
// ihour, imin: 时分
// dsec: 秒（十进制）
// gregflag: 历法标志
// 返回：ET儒略日, UT1儒略日, 错误信息
func UtcToJd(iyear, imonth, iday, ihour, imin Int32, dsec Float64, gregflag Int32) (et, ut1 Float64, err error) {
	if iyear < -4713 || iyear > 5000000 {
		return 0, 0, fmt.Errorf("年份超出有效范围: %d", iyear)
	}
	
	if imonth < 1 || imonth > 12 {
		return 0, 0, fmt.Errorf("无效月份: %d", imonth)
	}
	
	if iday < 1 || iday > 31 {
		return 0, 0, fmt.Errorf("无效日期: %d", iday)
	}
	
	if ihour < 0 || ihour >= 24 {
		return 0, 0, fmt.Errorf("无效小时: %d", ihour)
	}
	
	if imin < 0 || imin >= 60 {
		return 0, 0, fmt.Errorf("无效分钟: %d", imin)
	}
	
	if dsec < 0 || dsec >= 60 {
		return 0, 0, fmt.Errorf("无效秒数: %f", dsec)
	}
	
	// 计算总小时数
	dhour := Float64(ihour) + Float64(imin)/60.0 + dsec/3600.0
	
	// 计算UT1儒略日
	ut1 = Julday(int(iyear), int(imonth), int(iday), dhour, int(gregflag))
	
	// 计算Delta T并转换为ET
	dt := Deltat(ut1)
	et = ut1 + dt/86400.0
	
	return et, ut1, nil
}

// JdetToUtc 将ET儒略日转换为UTC时间
func JdetToUtc(tjdEt Float64, gregflag Int32) (iyear, imonth, iday, ihour, imin Int32, dsec Float64) {
	// 将ET转换为UT1
	dt := Deltat(tjdEt)
	tjdUt := tjdEt - dt/86400.0
	
	// 转换为日历日期
	year, month, day, jut := Revjul(tjdUt, int(gregflag))
	
	iyear = Int32(year)
	imonth = Int32(month)
	iday = Int32(day)
	
	// 转换小时分秒
	ihour = Int32(jut)
	remainder := jut - Float64(ihour)
	imin = Int32(remainder * 60.0)
	dsec = (remainder*60.0 - Float64(imin)) * 60.0
	
	return
}

// Jdut1ToUtc 将UT1儒略日转换为UTC时间
func Jdut1ToUtc(tjdUt Float64, gregflag Int32) (iyear, imonth, iday, ihour, imin Int32, dsec Float64) {
	// 转换为日历日期
	year, month, day, jut := Revjul(tjdUt, int(gregflag))
	
	iyear = Int32(year)
	imonth = Int32(month)
	iday = Int32(day)
	
	// 转换小时分秒
	ihour = Int32(jut)
	remainder := jut - Float64(ihour)
	imin = Int32(remainder * 60.0)
	dsec = (remainder*60.0 - Float64(imin)) * 60.0
	
	return
}

// UtcTimeZone 时区转换
func UtcTimeZone(iyear, imonth, iday, ihour, imin Int32, dsec, dTimezone Float64) (iyearOut, imonthOut, idayOut, ihourOut, iminOut Int32, dsecOut Float64) {
	// 计算总秒数
	totalSec := Float64(ihour)*3600 + Float64(imin)*60 + dsec
	
	// 应用时区偏移
	totalSec += dTimezone * 3600
	
	// 计算天数偏移
	dayOffset := int(math.Floor(totalSec / 86400))
	totalSec = math.Mod(totalSec, 86400)
	
	if totalSec < 0 {
		totalSec += 86400
		dayOffset--
	}
	
	// 转换为儒略日并加上偏移
	jd := Julday(int(iyear), int(imonth), int(iday), 0, SeGregCal)
	jd += Float64(dayOffset)
	
	// 转换回日历日期
	year, month, day, _ := Revjul(jd, SeGregCal)
	
	iyearOut = Int32(year)
	imonthOut = Int32(month)
	idayOut = Int32(day)
	
	// 转换时分秒
	ihourOut = Int32(totalSec / 3600)
	remainder := math.Mod(totalSec, 3600)
	iminOut = Int32(remainder / 60)
	dsecOut = math.Mod(remainder, 60)
	
	return
}

// DayOfWeek 计算星期几（0=周日，1=周一，...，6=周六）
func DayOfWeek(jd Float64) int {
	return int(math.Mod(math.Floor(jd+1.5), 7))
}

// 简化的Delta T计算（实际实现应该更复杂）
func Deltat(tjd Float64) Float64 {
	// 这里使用简化算法，实际应该根据历史数据和预测模型
	year := (tjd - J2000) / 365.25 + 2000.0
	
	if year < 1620 {
		// 1620年之前的估算
		t := (year - 2000) / 100
		return -20 + 32*t*t
	} else if year < 1972 {
		// 1620-1972年的估算
		t := year - 2000
		return 63.86 + 0.3345*t - 0.060374*t*t + 0.0017275*t*t*t + 0.000651814*t*t*t*t + 0.00002373599*t*t*t*t*t
	} else if year < 2005 {
		// 1972-2005年的实测值区间
		t := year - 2000
		return 63.86 + 0.3345*t - 0.060374*t*t + 0.0017275*t*t*t + 0.000651814*t*t*t*t + 0.00002373599*t*t*t*t*t
	} else {
		// 2005年之后的预测
		t := year - 2005
		return 64.184 + 0.8*t
	}
}

// GetCurrentTime 获取当前时间的儒略日
func GetCurrentTime() Float64 {
	now := time.Now().UTC()
	return Julday(now.Year(), int(now.Month()), now.Day(), 
		Float64(now.Hour()) + Float64(now.Minute())/60.0 + Float64(now.Second())/3600.0 + Float64(now.Nanosecond())/3600000000000.0, 
		SeGregCal)
}

// IsLeapYear 判断是否为闰年
func IsLeapYear(year int) bool {
	if year%4 != 0 {
		return false
	}
	if year%100 != 0 {
		return true
	}
	return year%400 == 0
}

// DaysInMonth 获取某月的天数
func DaysInMonth(year, month int) int {
	if month < 1 || month > 12 {
		return 0
	}
	
	days := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	
	if month == 2 && IsLeapYear(year) {
		return 29
	}
	
	return days[month-1]
}