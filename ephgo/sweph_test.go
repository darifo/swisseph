package ephgo

import (
	"math"
	"testing"
)

func TestConstants(t *testing.T) {
	// 测试基本常量
	if SeVersion != "2.10.03" {
		t.Errorf("SeVersion = %s, want 2.10.03", SeVersion)
	}
	
	if J2000 != 2451545.0 {
		t.Errorf("J2000 = %f, want 2451545.0", J2000)
	}
	
	if math.Abs(DegToRad - math.Pi/180.0) > 1e-15 {
		t.Errorf("DegToRad = %f, want %f", DegToRad, math.Pi/180.0)
	}
}

func TestJulday(t *testing.T) {
	// 测试儒略日计算
	tests := []struct {
		year, month, day int
		hour             Float64
		gregflag         int
		expected         Float64
	}{
		{2000, 1, 1, 12.0, SeGregCal, 2451545.0},    // J2000
		{1950, 1, 0, 22.0908333, SeGregCal, B1950},  // B1950（近似）
		{2024, 1, 1, 0.0, SeGregCal, 2460310.5},     // 2024年1月1日
	}
	
	for _, test := range tests {
		result := Julday(test.year, test.month, test.day, test.hour, test.gregflag)
		if math.Abs(result-test.expected) > 0.1 {
			t.Errorf("Julday(%d, %d, %d, %f, %d) = %f, want %f",
				test.year, test.month, test.day, test.hour, test.gregflag,
				result, test.expected)
		}
	}
}

func TestRevjul(t *testing.T) {
	// 测试儒略日反转
	jd := Float64(2460310.5) // 2024年1月1日
	year, month, day, jut := Revjul(jd, SeGregCal)
	
	if year != 2024 || month != 1 || day != 1 {
		t.Errorf("Revjul(%f) = %d-%d-%d, want 2024-1-1", jd, year, month, day)
	}
	
	if math.Abs(jut) > 0.001 {
		t.Errorf("Revjul(%f) hour = %f, want 0.0", jd, jut)
	}
}

func TestDateConversion(t *testing.T) {
	// 测试日期转换
	jd, err := DateConversion(2024, 1, 1, 12.0, 'g')
	if err != nil {
		t.Errorf("DateConversion failed: %v", err)
	}
	
	expected := Float64(2460311.0) // 2024年1月1日12:00
	if math.Abs(jd-expected) > 0.001 {
		t.Errorf("DateConversion(2024, 1, 1, 12.0, 'g') = %f, want %f", jd, expected)
	}
}

func TestDayOfWeek(t *testing.T) {
	// 测试星期几计算
	// 2024年1月1日是星期一
	jd := Julday(2024, 1, 1, 0.0, SeGregCal)
	dow := DayOfWeek(jd)
	
	if dow != 1 { // 1 = 星期一
		t.Errorf("DayOfWeek(2024-01-01) = %d, want 1 (Monday)", dow)
	}
}

func TestIsLeapYear(t *testing.T) {
	// 测试闰年判断
	tests := []struct {
		year     int
		expected bool
	}{
		{2000, true},  // 能被400整除
		{1900, false}, // 能被100整除但不能被400整除
		{2004, true},  // 能被4整除但不能被100整除
		{2003, false}, // 不能被4整除
		{2024, true},  // 闰年
	}
	
	for _, test := range tests {
		result := IsLeapYear(test.year)
		if result != test.expected {
			t.Errorf("IsLeapYear(%d) = %t, want %t", test.year, result, test.expected)
		}
	}
}

func TestDaysInMonth(t *testing.T) {
	// 测试月份天数
	tests := []struct {
		year, month int
		expected    int
	}{
		{2024, 2, 29}, // 闰年2月
		{2023, 2, 28}, // 平年2月
		{2024, 1, 31}, // 1月
		{2024, 4, 30}, // 4月
		{2024, 13, 0}, // 无效月份
	}
	
	for _, test := range tests {
		result := DaysInMonth(test.year, test.month)
		if result != test.expected {
			t.Errorf("DaysInMonth(%d, %d) = %d, want %d",
				test.year, test.month, result, test.expected)
		}
	}
}

func TestUtcToJd(t *testing.T) {
	// 测试UTC到儒略日转换
	et, ut1, err := UtcToJd(2024, 1, 1, 12, 0, 0.0, SeGregCal)
	if err != nil {
		t.Errorf("UtcToJd failed: %v", err)
	}
	
	// ET应该比UT1大（因为Delta T为正）
	if et <= ut1 {
		t.Errorf("ET (%f) should be greater than UT1 (%f)", et, ut1)
	}
	
	// Delta T应该在合理范围内（几十秒）
	deltaT := (et - ut1) * 86400
	if deltaT < 30 || deltaT > 100 {
		t.Errorf("Delta T = %f seconds, expected between 30-100", deltaT)
	}
}

func TestGetPlanetName(t *testing.T) {
	// 测试天体名称获取
	tests := []struct {
		ipl      int
		expected string
	}{
		{SeSun, SeNameSun},
		{SeMoon, SeNameMoon},
		{SeMercury, SeNameMercury},
		{SePluto, SeNamePluto},
		{SeMeanNode, SeNameMeanNode},
		{9999, "Planet_9999"}, // 未知天体
	}
	
	for _, test := range tests {
		result := GetPlanetName(test.ipl)
		if result != test.expected {
			t.Errorf("GetPlanetName(%d) = %s, want %s", test.ipl, result, test.expected)
		}
	}
}

func TestPolarToCartesian(t *testing.T) {
	// 测试极坐标到笛卡尔坐标转换
	polar := [6]Float64{
		0.0,  // 经度 0°
		0.0,  // 纬度 0°
		1.0,  // 距离 1 AU
		0.0,  // 经度速度
		0.0,  // 纬度速度
		0.0,  // 距离速度
	}
	
	cart := polarToCartesian(polar)
	
	// 在赤道面上，经度0°时，x=1, y=0, z=0
	if math.Abs(cart[0] - 1.0) > 1e-10 {
		t.Errorf("Cartesian X = %f, want 1.0", cart[0])
	}
	if math.Abs(cart[1]) > 1e-10 {
		t.Errorf("Cartesian Y = %f, want 0.0", cart[1])
	}
	if math.Abs(cart[2]) > 1e-10 {
		t.Errorf("Cartesian Z = %f, want 0.0", cart[2])
	}
}

func TestNormalizeAngle(t *testing.T) {
	// 测试角度归一化
	tests := []struct {
		input    Float64
		expected Float64
	}{
		{0.0, 0.0},
		{math.Pi, math.Pi},
		{2*math.Pi, 0.0},
		{3*math.Pi, math.Pi},
		{-math.Pi, math.Pi},
		{-2*math.Pi, 0.0},
	}
	
	for _, test := range tests {
		result := normalizeAngle(test.input)
		if math.Abs(result-test.expected) > 1e-10 {
			t.Errorf("normalizeAngle(%f) = %f, want %f", test.input, result, test.expected)
		}
	}
}

func TestVersion(t *testing.T) {
	// 测试版本获取
	version := Version()
	if version != SeVersion {
		t.Errorf("Version() = %s, want %s", version, SeVersion)
	}
}

func TestDegRadConversion(t *testing.T) {
	// 测试度和弧度转换
	deg := Float64(180.0)
	rad := deg2rad(deg)
	if math.Abs(rad - math.Pi) > 1e-10 {
		t.Errorf("deg2rad(180) = %f, want %f", rad, math.Pi)
	}
	
	back := rad2deg(rad)
	if math.Abs(back - deg) > 1e-10 {
		t.Errorf("rad2deg(π) = %f, want 180.0", back)
	}
}

func TestSquareSumAndDotProd(t *testing.T) {
	// 测试向量运算
	v1 := [3]Float64{3.0, 4.0, 0.0}
	v2 := [3]Float64{1.0, 0.0, 0.0}
	
	// 测试向量模长平方
	sum := SquareSum(v1)
	expected := Float64(25.0) // 3²+4²+0² = 25
	if math.Abs(sum-expected) > 1e-10 {
		t.Errorf("SquareSum([3,4,0]) = %f, want %f", sum, expected)
	}
	
	// 测试点积
	dot := DotProd(v1, v2)
	expected = Float64(3.0) // 3*1 + 4*0 + 0*0 = 3
	if math.Abs(dot-expected) > 1e-10 {
		t.Errorf("DotProd([3,4,0], [1,0,0]) = %f, want %f", dot, expected)
	}
}

// 基准测试
func BenchmarkJulday(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Julday(2024, 1, 1, 12.0, SeGregCal)
	}
}

func BenchmarkRevjul(b *testing.B) {
	jd := Float64(2460310.5)
	for i := 0; i < b.N; i++ {
		Revjul(jd, SeGregCal)
	}
}

func BenchmarkDeltat(b *testing.B) {
	jd := Float64(2460310.5)
	for i := 0; i < b.N; i++ {
		Deltat(jd)
	}
}