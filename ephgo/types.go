package ephgo

import (
	"os"
	"sync"
)

// Bool 类型定义
type Bool bool

// 基础类型别名
type (
	Int32    = int32
	Uint32   = uint32
	Int16    = int16
	Float64  = float64
	Centisec = int32  // 厘秒用于角度和时间
)

// FileData 星历文件数据结构
type FileData struct {
	Fnam        string  // 星历文件名
	Fversion    int     // 文件版本号
	Astnam      string  // 小行星名称
	SwephDenum  Int32   // JPL星历DE编号
	Fptr        *os.File // 文件指针
	Tfstart     Float64 // 文件有效开始时间
	Tfend       Float64 // 文件有效结束时间
	Iflg        Int32   // 字节顺序和大小端标志
	Npl         int16   // 文件中行星数量
	Ipl         [SeiFileNmaxplan]int // 行星编号数组
}

// PlanData 行星数据结构
type PlanData struct {
	// 文件打开时读取的数据
	Ibdy   int     // 内部天体编号
	Iflg   Int32   // 标志位
	Ncoe   int     // 星历多项式系数数量
	Lndx0  Int32   // 行星索引在文件中的位置
	Nndx   Int32   // 文件中索引条目数量
	Tfstart Float64 // 文件包含的星历起始时间
	Tfend   Float64 // 文件包含的星历结束时间
	Dseg    Float64 // 段大小（多项式覆盖的天数）
	
	// 轨道要素
	Telem Float64 // 要素历元
	Prot  Float64
	Qrot  Float64
	Dprot Float64
	Dqrot Float64
	Rmax  Float64 // 切比雪夫系数的归一化因子
	
	// 参考椭圆（如果使用）
	Peri  Float64
	Dperi Float64
	Refep []Float64 // 参考椭圆切比雪夫系数指针
	
	// 当前段信息
	Tseg0 Float64   // 当前段起始JD
	Tseg1 Float64   // 当前段结束JD
	Segp  []Float64 // 当前段切比雪夫系数
	Neval int       // 需要计算的系数数量
	
	// 最近计算结果
	Teval   Float64     // 上次计算的时间
	Iephe   Int32       // 使用的星历类型
	X       [6]Float64  // 位置和速度矢量（J2000赤道坐标）
	Xflgs   Int32       // 标志位
	Xreturn [24]Float64 // 返回位置：极坐标、笛卡尔坐标等
}

// NodeData 月球交点和远地点数据
type NodeData struct {
	Teval   Float64     // 上次计算时间
	Iephe   Int32       // 使用的星历
	X       [6]Float64  // 位置和速度矢量
	Xflgs   Int32       // 标志位
	Xreturn [24]Float64 // 返回位置
}

// Epsilon 黄赤交角数据
type Epsilon struct {
	Teps Float64 // 儒略日
	Eps  Float64 // 黄赤交角
	Seps Float64 // sin(黄赤交角)
	Ceps Float64 // cos(黄赤交角)
}

// Nut 章动数据
type Nut struct {
	Tnut   Float64      // 计算章动的时间
	Nutlo  [2]Float64   // 经度和倾角章动
	Snut   Float64      // sin(倾角章动)
	Cnut   Float64      // cos(倾角章动)
	Matrix [3][3]Float64 // 章动矩阵
}

// GenConst 一般常数
type GenConst struct {
	Clight       Float64 // 光速
	Aunit        Float64 // 天文单位
	Helgravconst Float64 // 太阳引力常数
	Ratme        Float64 // 地月质量比
	Sunradius    Float64 // 太阳半径
}

// SavePositions 保存的位置数据
type SavePositions struct {
	Ipl      int         // 行星编号
	Tsave    Float64     // 保存时间
	Iflgsave Int32       // 保存的标志
	Xsaves   [24]Float64 // 保存的位置坐标
}

// TopoData 地心数据
type TopoData struct {
	Geolon Float64     // 地理经度
	Geolat Float64     // 地理纬度
	Geoalt Float64     // 地理高度
	Teval  Float64     // 计算时间
	TjdUt  Float64     // UT时间
	Xobs   [6]Float64  // 观测者位置
}

// SidData 恒星时数据
type SidData struct {
	SidMode Int32   // 恒星时模式
	AyanT0  Float64 // 起始岁差
	T0      Float64 // 参考时间
	T0IsUT  Bool    // T0是否为UT
}

// FixedStar 恒星数据
type FixedStar struct {
	Skey     string  // 搜索关键字
	Starname string  // 恒星名称
	Starbayer string // 拜耳名称
	Starno   string  // 星表号
	Epoch    Float64 // 历元
	Ra       Float64 // 赤经
	De       Float64 // 赤纬
	Ramot    Float64 // 赤经自行
	Demot    Float64 // 赤纬自行
	Radvel   Float64 // 径向速度
	Parall   Float64 // 视差
	Mag      Float64 // 星等
}

// Interpol 插值数据
type Interpol struct {
	TjdNut0   Float64 // 章动计算时间0
	TjdNut2   Float64 // 章动计算时间2
	NutDpsi0  Float64 // 章动经度0
	NutDpsi1  Float64 // 章动经度1
	NutDpsi2  Float64 // 章动经度2
	NutDeps0  Float64 // 章动倾角0
	NutDeps1  Float64 // 章动倾角1
	NutDeps2  Float64 // 章动倾角2
}

// AyaInit 岁差初始化数据
type AyaInit struct {
	T0        Float64 // 岁差历元
	AyanT0    Float64 // 岁差值
	T0IsUT    Bool    // T0是否为UT
	PrecOffset int    // 岁差偏移
}

// PlantBl 行星表数据
type PlantBl struct {
	MaxHarmonic [9]uint8      // 最大谐波
	MaxPowerOfT uint8         // t的最大幂次
	ArgTbl      []int8        // 参数表
	LonTbl      []Float64     // 经度表
	LatTbl      []Float64     // 纬度表
	RadTbl      []Float64     // 半径表
	Distance    Float64       // 距离
}

// SweData 主要数据结构（全局状态）
type SweData struct {
	EphePathIsSet         Bool      // 星历路径是否已设置
	JplFileIsOpen         Bool      // JPL文件是否已打开
	Fixfp                 *os.File  // 恒星文件指针
	Ephepath              string    // 星历路径
	Jplfnam               string    // JPL文件名
	Jpldenum              Int32     // JPL DE编号
	LastEpheflag          Int32     // 最后使用的星历标志
	GeoposIsSet           Bool      // 地理位置是否已设置
	AyanaIsSet            Bool      // 岁差是否已设置
	IsOldStarfile         Bool      // 是否为旧恒星文件
	EopTjdBeg             Float64   // EOP起始时间
	EopTjdBegHorizons     Float64   // Horizons EOP起始时间
	EopTjdEnd             Float64   // EOP结束时间
	EopTjdEndAdd          Float64   // EOP附加结束时间
	EopDpsiLoaded         int       // EOP章动是否已加载
	TidAcc                Float64   // 潮汐加速度
	IsTidAccManual        Bool      // 是否手动设置潮汐加速度
	InitDtDone            Bool      // DT初始化是否完成
	SwedIsInitialised     Bool      // 是否已初始化
	DeltaTUserdefIsSet    Bool      // 用户定义DT是否已设置
	DeltaTUserdef         Float64   // 用户定义DT
	AstG                  Float64   // 小行星G参数
	AstH                  Float64   // 小行星H参数
	AstDiam               Float64   // 小行星直径
	Astelem               string    // 小行星轨道要素
	ISavedPlanetName      int       // 保存的行星名称索引
	SavedPlanetName       string    // 保存的行星名称
	Dpsi                  []Float64 // 章动经度数组
	Deps                  []Float64 // 章动倾角数组
	Timeout               Int32     // 超时
	AstroModels           [SeiNmodels]Int32 // 天体模型
	DoInterpolateNut      Bool      // 是否插值章动
	Interpol              Interpol  // 插值数据
	Fidat                 [SeiNephfiles]FileData // 文件数据
	Gcdat                 GenConst  // 一般常数
	Pldat                 [SeiNplanets]PlanData // 行星数据
	Nddat                 [SeiNnodeEtc]PlanData // 交点数据
	Savedat               [SeNplanets + 1]SavePositions // 保存位置
	Oec                   Epsilon   // 真黄赤交角
	Oec2000               Epsilon   // J2000黄赤交角
	Nut                   Nut       // 章动
	Nut2000               Nut       // J2000章动
	Nutv                  Nut       // 速度章动
	Topd                  TopoData  // 地心数据
	Sidd                  SidData   // 恒星时数据
	NFixstarsReal         Bool      // 实际恒星数量
	NFixstarsNamed        Bool      // 命名恒星数量
	NFixstarsRecords      Bool      // 恒星记录数量
	FixedStars            []FixedStar // 恒星数组
}

// 全局数据实例（线程安全）
var (
	swed     SweData
	swedMu   sync.RWMutex
)

// GetSweData 获取全局数据（线程安全）
func GetSweData() *SweData {
	swedMu.RLock()
	defer swedMu.RUnlock()
	return &swed
}

// SetSweData 设置全局数据（线程安全）
func SetSweData(data *SweData) {
	swedMu.Lock()
	defer swedMu.Unlock()
	swed = *data
}

// 数学辅助函数
func SquareSum(x [3]Float64) Float64 {
	return x[0]*x[0] + x[1]*x[1] + x[2]*x[2]
}

func DotProd(x, y [3]Float64) Float64 {
	return x[0]*y[0] + x[1]*y[1] + x[2]*y[2]
}