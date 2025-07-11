package ephgo

import "math"

// 版本信息
const SeVersion = "2.10.03"

// 基础天文常量
const (
	J2000 = 2451545.0            // 2000 January 1.5
	B1950 = 2433282.42345905     // 1950 January 0.923
	J1900 = 2415020.0            // 1900 January 0.5
	B1850 = 2396758.2035810      // 1850 January 16:53

	Pi     = math.Pi
	TwoPi  = 2.0 * math.Pi
	DegToRad = math.Pi / 180.0
	RadToDeg = 180.0 / math.Pi
)

// 小行星编号
const (
	MpcCeres  = 1
	MpcPallas = 2
	MpcJuno   = 3
	MpcVesta  = 4
	MpcChiron = 2060
	MpcPholus = 5145
)

// 天体名称
const (
	SeNameSun         = "Sun"
	SeNameMoon        = "Moon" 
	SeNameMercury     = "Mercury"
	SeNameVenus       = "Venus"
	SeNameMars        = "Mars"
	SeNameJupiter     = "Jupiter"
	SeNameSaturn      = "Saturn"
	SeNameUranus      = "Uranus"
	SeNameNeptune     = "Neptune"
	SeNamePluto       = "Pluto"
	SeNameMeanNode    = "mean Node"
	SeNameTrueNode    = "true Node"
	SeNameMeanApog    = "mean Apogee"
	SeNameOscuApog    = "osc. Apogee"
	SeNameIntpApog    = "intp. Apogee"
	SeNameIntpPerg    = "intp. Perigee"
	SeNameEarth       = "Earth"
	SeNameCeres       = "Ceres"
	SeNamePallas      = "Pallas"
	SeNameJuno        = "Juno"
	SeNameVesta       = "Vesta"
	SeNameChiron      = "Chiron"
	SeNamePholus      = "Pholus"
)

// 天体编号（用于计算函数的ipl参数）
const (
	SeEclNut = -1

	SeSun      = 0
	SeMoon     = 1
	SeMercury  = 2
	SeVenus    = 3
	SeMars     = 4
	SeJupiter  = 5
	SeSaturn   = 6
	SeUranus   = 7
	SeNeptune  = 8
	SePluto    = 9
	SeMeanNode = 10
	SeTrueNode = 11
	SeMeanApog = 12
	SeOscuApog = 13
	SeEarth    = 14
	SeChiron   = 15
	SePholus   = 16
	SeCeres    = 17
	SePallas   = 18
	SeJuno     = 19
	SeVesta    = 20
	SeIntpApog = 21
	SeIntpPerg = 22

	SeNplanets = 23
)

// 内部天体索引
const (
	SeiEpsilon = -2
	SeiNutation = -1
	SeiEmb      = 0  // Earth-Moon barycenter
	SeiEarth    = 0
	SeiSun      = 0
	SeiMoon     = 1
	SeiMercury  = 2
	SeiVenus    = 3
	SeiMars     = 4
	SeiJupiter  = 5
	SeiSaturn   = 6
	SeiUranus   = 7
	SeiNeptune  = 8
	SeiPluto    = 9
	SeiSunbary  = 10 // barycentric sun
	SeiAnybody  = 11 // any asteroid
	SeiChiron   = 12
	SeiPholus   = 13
	SeiCeres    = 14
	SeiPallas   = 15
	SeiJuno     = 16
	SeiVesta    = 17

	SeiNplanets = 18

	SeiMeanNode = 0
	SeiTrueNode = 1
	SeiMeanApog = 2
	SeiOscuApog = 3
	SeiIntpApog = 4
	SeiIntpPerg = 5

	SeiNnodeEtc = 6
)

// 计算标志位
const (
	SeflgJpleph   = 1       // 使用JPL星历
	SeflgSwieph   = 2       // 使用SWISS EPHEMERIS星历
	SeflgMoseph   = 4       // 使用Moshier星历
	
	SeflgHelctr   = 8       // 日心坐标
	SeflgTruepos  = 16      // 真实/几何位置，非视位置
	SeflgJ2000    = 32      // 无岁差，即J2000春分点
	SeflgNonut    = 64      // 无章动，即平春分点
	SeflgSpeed3   = 128     // 3点法计算速度（不推荐使用）
	SeflgSpeed    = 256     // 高精度速度
	SeflgNogdefl  = 512     // 关闭引力偏折
	SeflgNoaberr  = 1024    // 关闭光行差
	SeflgAstrometric = SeflgNoaberr | SeflgNogdefl // 天体测量位置
	
	SeflgEquatorial = 2048  // 赤道坐标
	SeflgXyz        = 4096  // 笛卡尔坐标
	SeflgRadians    = 8192  // 弧度而非度数
	SeflgBaryctr    = 16384 // 重心坐标
	SeflgTopoctr    = 32768 // 地心坐标
	
	SeflgTropical   = 0     // 回归坐标（默认）
	SeflgSidereal   = 65536 // 恒星坐标
	SeflgIcrs       = 131072 // ICRS参考系
)

// 日历类型
const (
	SeJulCal  = 0 // 儒略历
	SeGregCal = 1 // 格里高利历
)

// 物理常数
const (
	// 月球参数
	MoonMeanDist = 384400000.0  // 平均距离（米）
	MoonMeanIncl = 5.1453964    // 平均倾角
	MoonMeanEcc  = 0.054900489  // 平均偏心率
	
	// 地球参数
	SunEarthMrat    = 332946.050895    // 太阳质量/地球质量
	EarthMoonMrat   = 1 / 0.0123000383 // 地球质量/月球质量
	EarthRadius     = 6378136.6        // 地球半径（米）
	EarthOblateness = 1.0 / 298.25642  // 地球扁率
	EarthRotSpeed   = 7.2921151467e-5 * 86400 // 地球自转速度（弧度/天）
	
	// 基本物理常数
	Aunit        = 1.49597870700e+11 // 天文单位（米）
	Clight       = 2.99792458e+8     // 光速（米/秒）
	Helgravconst = 1.32712440017987e+20 // 太阳引力常数
	Geogconst    = 3.98600448e+14       // 地球引力常数
	Kgauss       = 0.01720209895        // 高斯引力常数
	
	// 光行时和视差
	LighttimeAunit = 499.0047838362 / 3600.0 / 24.0 // 8.3167分钟（天）
	ParsecToAunit  = 206264.8062471                 // 秒差距到天文单位
)

// 文件常量
const (
	SeiFileNmaxplan   = 50
	SeiFileEfposbegin = 500
	SeFileSuffix      = "se1"
	SeiNephfiles      = 7
	SeiCurrFpos       = -1
	SeiNmodels        = 8
)

// JPL天体索引
const (
	JMercury = 0
	JVenus   = 1
	JEarth   = 2
	JMars    = 3
	JJupiter = 4
	JSaturn  = 5
	JUranus  = 6
	JNeptune = 7
	JPluto   = 8
	JMoon    = 9
	JSun     = 10
	JSbary   = 11
	JEmb     = 12
	JNut     = 13
	JLib     = 14
)

// 最大字符串长度
const AsMaxch = 256

// 返回状态码
const (
	Ok                = 0
	Err               = -1
	NotAvailable      = -2
	BeyondEphLimits   = -3
)

// 天体直径（米）
var PlaDiam = [SeVesta + 1]float64{
	1392000000.0, // Sun
	3475000.0,    // Moon
	2439400.0 * 2, // Mercury
	6051800.0 * 2, // Venus
	3389500.0 * 2, // Mars
	69911000.0 * 2, // Jupiter
	58232000.0 * 2, // Saturn
	25362000.0 * 2, // Uranus
	24622000.0 * 2, // Neptune
	1188300.0 * 2, // Pluto
	0, 0, 0, 0,    // nodes and apogees
	6371008.4 * 2, // Earth
	271370.0,      // Chiron
	290000.0,      // Pholus
	939400.0,      // Ceres
	545000.0,      // Pallas
	246596.0,      // Juno
	525400.0,      // Vesta
}
