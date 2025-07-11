package ephgo

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
)

// 全局变量
var (
	ephePath     string = ""
	jplFileName  string = ""
	isInitialized bool  = false
)

// SetEphePath 设置星历文件路径
func SetEphePath(path string) {
	ephePath = path
	swed := GetSweData()
	swed.EphePathIsSet = true
	swed.Ephepath = path
	SetSweData(swed)
}

// SetJplFile 设置JPL文件名
func SetJplFile(fname string) {
	jplFileName = fname
	swed := GetSweData()
	swed.Jplfnam = fname
	SetSweData(swed)
}

// Close 关闭Swiss Ephemeris并释放资源
func Close() {
	CloseJplFile()
	swed := GetSweData()
	
	// 关闭文件
	for i := range swed.Fidat {
		if swed.Fidat[i].Fptr != nil {
			swed.Fidat[i].Fptr.Close()
			swed.Fidat[i].Fptr = nil
		}
	}
	
	if swed.Fixfp != nil {
		swed.Fixfp.Close()
		swed.Fixfp = nil
	}
	
	// 重置状态
	swed.SwedIsInitialised = false
	swed.EphePathIsSet = false
	swed.JplFileIsOpen = false
	
	SetSweData(swed)
	isInitialized = false
}

// Calc 计算天体位置
// tjd: 儒略日（力学时ET）
// ipl: 天体编号
// iflag: 计算标志
// 返回：坐标数组xx[6]，错误信息
func Calc(tjd Float64, ipl int, iflag Int32) ([6]Float64, error) {
	var xx [6]Float64
	
	// 检查输入参数
	if tjd < -5000000 || tjd > 5000000 {
		return xx, fmt.Errorf("儒略日超出有效范围: %f", tjd)
	}
	
	// 初始化
	err := initializeSwissEph()
	if err != nil {
		return xx, fmt.Errorf("初始化失败: %v", err)
	}
	
	// 根据天体类型分发计算
	switch {
	case ipl >= SeSun && ipl <= SePluto:
		return calcMainPlanet(tjd, ipl, iflag)
	case ipl == SeMeanNode || ipl == SeTrueNode:
		return calcNode(tjd, ipl, iflag)
	case ipl == SeMeanApog || ipl == SeOscuApog:
		return calcApogee(tjd, ipl, iflag)
	case ipl >= SeCeres && ipl <= SeVesta:
		return calcAsteroid(tjd, ipl, iflag)
	case ipl >= SeChiron && ipl <= SePholus:
		return calcMinorPlanet(tjd, ipl, iflag)
	default:
		return xx, fmt.Errorf("不支持的天体编号: %d", ipl)
	}
}

// CalcUT 计算天体位置（世界时UT）
func CalcUT(tjdUt Float64, ipl int, iflag Int32) ([6]Float64, error) {
	// 转换UT到ET
	dt := Deltat(tjdUt)
	tjdEt := tjdUt + dt/86400.0
	
	return Calc(tjdEt, ipl, iflag)
}

// Version 获取版本信息
func Version() string {
	return SeVersion
}

// GetPlanetName 获取天体名称
func GetPlanetName(ipl int) string {
	switch ipl {
	case SeSun:
		return SeNameSun
	case SeMoon:
		return SeNameMoon
	case SeMercury:
		return SeNameMercury
	case SeVenus:
		return SeNameVenus
	case SeMars:
		return SeNameMars
	case SeJupiter:
		return SeNameJupiter
	case SeSaturn:
		return SeNameSaturn
	case SeUranus:
		return SeNameUranus
	case SeNeptune:
		return SeNameNeptune
	case SePluto:
		return SeNamePluto
	case SeMeanNode:
		return SeNameMeanNode
	case SeTrueNode:
		return SeNameTrueNode
	case SeMeanApog:
		return SeNameMeanApog
	case SeOscuApog:
		return SeNameOscuApog
	case SeEarth:
		return SeNameEarth
	case SeCeres:
		return SeNameCeres
	case SePallas:
		return SeNamePallas
	case SeJuno:
		return SeNameJuno
	case SeVesta:
		return SeNameVesta
	case SeChiron:
		return SeNameChiron
	case SePholus:
		return SeNamePholus
	default:
		return fmt.Sprintf("Planet_%d", ipl)
	}
}

// initializeSwissEph 初始化Swiss Ephemeris
func initializeSwissEph() error {
	if isInitialized {
		return nil
	}
	
	swed := GetSweData()
	
	// 设置默认路径
	if ephePath == "" {
		ephePath = getDefaultEphePath()
	}
	
	// 初始化常数
	swed.Gcdat.Clight = Clight
	swed.Gcdat.Aunit = Aunit
	swed.Gcdat.Helgravconst = Helgravconst
	swed.Gcdat.Ratme = EarthMoonMrat
	swed.Gcdat.Sunradius = PlaDiam[SeSun] / 2.0
	
	// 标记已初始化
	swed.SwedIsInitialised = true
	SetSweData(swed)
	isInitialized = true
	
	return nil
}

// getDefaultEphePath 获取默认星历路径
func getDefaultEphePath() string {
	// 检查环境变量
	if path := os.Getenv("SE_EPHE_PATH"); path != "" {
		return path
	}
	
	// 检查当前目录
	if _, err := os.Stat("./ephe"); err == nil {
		return "./ephe"
	}
	
	// 返回默认路径
	return "./"
}

// calcMainPlanet 计算主要行星位置
func calcMainPlanet(tjd Float64, ipl int, iflag Int32) ([6]Float64, error) {
	var xx [6]Float64
	
	// 优先使用JPL星历
	if (iflag&SeflgJpleph) != 0 || (iflag&SeflgSwieph) == 0 {
		if IsJplAvailable() {
			return calcWithJPL(tjd, ipl, iflag)
		}
	}
	
	// 使用Swiss Ephemeris文件
	if (iflag & SeflgSwieph) != 0 {
		return calcWithSwissEph(tjd, ipl, iflag)
	}
	
	// 使用Moshier星历
	if (iflag & SeflgMoseph) != 0 {
		return calcWithMoshier(tjd, ipl, iflag)
	}
	
	return xx, fmt.Errorf("没有可用的星历数据")
}

// calcWithJPL 使用JPL星历计算
func calcWithJPL(tjd Float64, ipl int, iflag Int32) ([6]Float64, error) {
	var xx [6]Float64
	
	// 映射天体编号到JPL编号
	jplTarget := mapToJPLTarget(ipl)
	jplCenter := JSbary // 默认相对于太阳系质心
	
	if (iflag & SeflgHelctr) != 0 {
		jplCenter = JSun // 日心坐标
	}
	
	// 调用JPL计算
	result, err := Pleph(tjd, jplTarget, jplCenter)
	if err != nil {
		return xx, err
	}
	
	// 复制结果
	copy(xx[:], result[:])
	
	// 应用坐标变换
	return applyCoordinateTransforms(xx, iflag), nil
}

// calcWithSwissEph 使用Swiss Ephemeris文件计算
func calcWithSwissEph(tjd Float64, ipl int, iflag Int32) ([6]Float64, error) {
	var xx [6]Float64
	
	// 这里应该实现Swiss Ephemeris文件读取和计算
	// 简化实现，实际需要读取.se1文件
	
	return xx, fmt.Errorf("Swiss Ephemeris文件计算尚未实现")
}

// calcWithMoshier 使用Moshier星历计算
func calcWithMoshier(tjd Float64, ipl int, iflag Int32) ([6]Float64, error) {
	var xx [6]Float64
	
	// 这里应该实现Moshier理论计算
	// 简化实现
	
	return xx, fmt.Errorf("Moshier星历计算尚未实现")
}

// calcNode 计算月球交点
func calcNode(tjd Float64, ipl int, iflag Int32) ([6]Float64, error) {
	var xx [6]Float64
	
	// 简化的月球交点计算
	// 实际实现需要更复杂的算法
	
	t := (tjd - J2000) / 36525.0
	
	if ipl == SeMeanNode {
		// 平交点
		lon := 125.0445479 - 1934.1362891*t + 0.0020754*t*t
		lat := 0.0
		rad := 1.0
		
		xx[0] = math.Mod(lon*DegToRad, TwoPi)
		xx[1] = lat * DegToRad
		xx[2] = rad
		
		// 速度（简化）
		xx[3] = -1934.1362891 * DegToRad / 36525.0
		xx[4] = 0.0
		xx[5] = 0.0
	} else {
		// 真交点（需要更复杂的计算）
		return calcNode(tjd, SeMeanNode, iflag)
	}
	
	return applyCoordinateTransforms(xx, iflag), nil
}

// calcApogee 计算月球远地点
func calcApogee(tjd Float64, ipl int, iflag Int32) ([6]Float64, error) {
	var xx [6]Float64
	
	// 简化的远地点计算
	t := (tjd - J2000) / 36525.0
	
	if ipl == SeMeanApog {
		// 平远地点
		lon := 83.3532465 + 4069.0137287*t - 0.0103200*t*t
		lat := 0.0
		rad := 1.0
		
		xx[0] = math.Mod(lon*DegToRad, TwoPi)
		xx[1] = lat * DegToRad
		xx[2] = rad
		
		// 速度
		xx[3] = 4069.0137287 * DegToRad / 36525.0
		xx[4] = 0.0
		xx[5] = 0.0
	} else {
		// 振荡远地点
		return calcApogee(tjd, SeMeanApog, iflag)
	}
	
	return applyCoordinateTransforms(xx, iflag), nil
}

// calcAsteroid 计算小行星位置
func calcAsteroid(tjd Float64, ipl int, iflag Int32) ([6]Float64, error) {
	var xx [6]Float64
	
	// 这里应该读取小行星轨道要素文件并计算
	// 简化实现
	
	return xx, fmt.Errorf("小行星计算尚未实现")
}

// calcMinorPlanet 计算小天体位置
func calcMinorPlanet(tjd Float64, ipl int, iflag Int32) ([6]Float64, error) {
	var xx [6]Float64
	
	// 这里应该计算Chiron、Pholus等小天体
	// 简化实现
	
	return xx, fmt.Errorf("小天体计算尚未实现")
}

// mapToJPLTarget 映射天体编号到JPL编号
func mapToJPLTarget(ipl int) int {
	switch ipl {
	case SeSun:
		return JSun
	case SeMoon:
		return JMoon
	case SeMercury:
		return JMercury
	case SeVenus:
		return JVenus
	case SeMars:
		return JMars
	case SeJupiter:
		return JJupiter
	case SeSaturn:
		return JSaturn
	case SeUranus:
		return JUranus
	case SeNeptune:
		return JNeptune
	case SePluto:
		return JPluto
	case SeEarth:
		return JEarth
	default:
		return JSun // 默认
	}
}

// applyCoordinateTransforms 应用坐标变换
func applyCoordinateTransforms(xx [6]Float64, iflag Int32) [6]Float64 {
	var result [6]Float64 = xx
	
	// 极坐标到笛卡尔坐标
	if (iflag & SeflgXyz) != 0 {
		result = polarToCartesian(result)
	}
	
	// 角度单位转换
	if (iflag & SeflgRadians) == 0 && (iflag & SeflgXyz) == 0 {
		// 转换为度
		result[0] *= RadToDeg
		result[1] *= RadToDeg
		result[3] *= RadToDeg
		result[4] *= RadToDeg
	}
	
	// 其他变换...
	
	return result
}

// polarToCartesian 极坐标转笛卡尔坐标
func polarToCartesian(polar [6]Float64) [6]Float64 {
	var cart [6]Float64
	
	lon := polar[0]
	lat := polar[1]
	rad := polar[2]
	dlon := polar[3]
	dlat := polar[4]
	drad := polar[5]
	
	cosLat := math.Cos(lat)
	sinLat := math.Sin(lat)
	cosLon := math.Cos(lon)
	sinLon := math.Sin(lon)
	
	// 位置
	cart[0] = rad * cosLat * cosLon
	cart[1] = rad * cosLat * sinLon
	cart[2] = rad * sinLat
	
	// 速度
	cart[3] = drad*cosLat*cosLon - rad*dlat*sinLat*cosLon - rad*dlon*cosLat*sinLon
	cart[4] = drad*cosLat*sinLon - rad*dlat*sinLat*sinLon + rad*dlon*cosLat*cosLon
	cart[5] = drad*sinLat + rad*dlat*cosLat
	
	return cart
}

// normalizeAngle 角度归一化到[0, 2π)
func normalizeAngle(angle Float64) Float64 {
	result := math.Mod(angle, TwoPi)
	if result < 0 {
		result += TwoPi
	}
	return result
}

// deg2rad 度转弧度
func deg2rad(deg Float64) Float64 {
	return deg * DegToRad
}

// rad2deg 弧度转度
func rad2deg(rad Float64) Float64 {
	return rad * RadToDeg
}

// findEphemerisFile 查找星历文件
func findEphemerisFile(filename string) (string, error) {
	// 检查完整路径
	if filepath.IsAbs(filename) {
		if _, err := os.Stat(filename); err == nil {
			return filename, nil
		}
		return "", fmt.Errorf("文件不存在: %s", filename)
	}
	
	// 在星历路径中查找
	searchPaths := []string{ephePath, "./ephe", "./"}
	
	for _, path := range searchPaths {
		if path == "" {
			continue
		}
		
		fullPath := filepath.Join(path, filename)
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath, nil
		}
	}
	
	return "", fmt.Errorf("找不到星历文件: %s", filename)
}

// 获取当前文件数据
func GetCurrentFileData(ifno int) (tfstart, tfend Float64, denum int) {
	swed := GetSweData()
	
	if ifno >= 0 && ifno < len(swed.Fidat) {
		return swed.Fidat[ifno].Tfstart, swed.Fidat[ifno].Tfend, int(swed.Fidat[ifno].SwephDenum)
	}
	
	return 0, 0, 0
}