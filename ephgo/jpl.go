package ephgo

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// JPL文件相关常量
const (
	JplRecordSize = 8144  // JPL记录大小
	JplNCoeff     = 1018  // 系数数量
	JplMaxPlanets = 15    // 最大行星数
)

// JplHeader JPL文件头结构
type JplHeader struct {
	StartJD  Float64   // 起始儒略日
	EndJD    Float64   // 结束儒略日
	StepJD   Float64   // 步长
	NumConst int       // 常数数量
	AU       Float64   // 天文单位
	EMRat    Float64   // 地月质量比
	IPT      [3][13]int // 指针表
	NumCoeff [13]int   // 系数数量
	Constants map[string]Float64 // 常数表
}

// JplData JPL数据结构
type JplData struct {
	Header   JplHeader
	File     *os.File
	FileName string
	IsOpen   bool
}

var jplData JplData

// OpenJplFile 打开JPL文件
func OpenJplFile(ss []Float64, fname, fpath string) error {
	if jplData.IsOpen {
		CloseJplFile()
	}
	
	// 构建文件路径
	var fullPath string
	if fpath != "" {
		fullPath = filepath.Join(fpath, fname)
	} else {
		fullPath = fname
	}
	
	// 打开文件
	file, err := os.Open(fullPath)
	if err != nil {
		return fmt.Errorf("无法打开JPL文件 %s: %v", fullPath, err)
	}
	
	jplData.File = file
	jplData.FileName = fullPath
	
	// 读取文件头
	err = readJplHeader(&jplData.Header, file)
	if err != nil {
		file.Close()
		return fmt.Errorf("读取JPL文件头失败: %v", err)
	}
	
	jplData.IsOpen = true
	
	// 返回开始、结束和步长
	if len(ss) >= 3 {
		ss[0] = jplData.Header.StartJD
		ss[1] = jplData.Header.EndJD
		ss[2] = jplData.Header.StepJD
	}
	
	return nil
}

// CloseJplFile 关闭JPL文件
func CloseJplFile() {
	if jplData.IsOpen && jplData.File != nil {
		jplData.File.Close()
		jplData.IsOpen = false
	}
}

// GetJplDenum 获取JPL DE编号
func GetJplDenum() Int32 {
	// 从文件名推断DE编号
	fname := strings.ToLower(filepath.Base(jplData.FileName))
	if strings.Contains(fname, "de405") {
		return 405
	} else if strings.Contains(fname, "de406") {
		return 406
	} else if strings.Contains(fname, "de431") {
		return 431
	} else if strings.Contains(fname, "de430") {
		return 430
	}
	return 431 // 默认
}

// Pleph 计算天体位置
// et: 儒略历力学时
// ntarg: 目标天体
// ncent: 中心天体
// 返回：位置和速度数组 [x, y, z, vx, vy, vz]
func Pleph(et Float64, ntarg, ncent int) ([6]Float64, error) {
	var rrd [6]Float64
	
	if !jplData.IsOpen {
		return rrd, fmt.Errorf("JPL文件未打开")
	}
	
	// 检查时间范围
	if et < jplData.Header.StartJD || et > jplData.Header.EndJD {
		return rrd, fmt.Errorf("时间超出JPL文件范围: %f", et)
	}
	
	// 计算记录位置
	recordNum := int((et - jplData.Header.StartJD) / jplData.Header.StepJD)
	if recordNum < 0 {
		recordNum = 0
	}
	
	// 读取记录
	record, err := readJplRecord(jplData.File, recordNum)
	if err != nil {
		return rrd, fmt.Errorf("读取JPL记录失败: %v", err)
	}
	
	// 计算位置和速度
	err = computePosition(record, et, ntarg, ncent, &rrd)
	if err != nil {
		return rrd, fmt.Errorf("计算位置失败: %v", err)
	}
	
	return rrd, nil
}

// readJplHeader 读取JPL文件头
func readJplHeader(header *JplHeader, file *os.File) error {
	// 定位到文件开始
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	
	// 读取第一个记录（文件头）
	buf := make([]byte, JplRecordSize)
	_, err = io.ReadFull(file, buf)
	if err != nil {
		return err
	}
	
	// 解析头部信息
	reader := NewByteReader(buf, binary.LittleEndian)
	
	// 读取基本参数
	header.StartJD = reader.ReadFloat64()
	header.EndJD = reader.ReadFloat64()
	header.StepJD = reader.ReadFloat64()
	header.NumConst = int(reader.ReadInt32())
	header.AU = reader.ReadFloat64()
	header.EMRat = reader.ReadFloat64()
	
	// 读取指针表
	for i := 0; i < 3; i++ {
		for j := 0; j < 13; j++ {
			header.IPT[i][j] = int(reader.ReadInt32())
		}
	}
	
	// 读取系数数量
	for i := 0; i < 13; i++ {
		header.NumCoeff[i] = int(reader.ReadInt32())
	}
	
	// 初始化常数表
	header.Constants = make(map[string]Float64)
	
	return nil
}

// readJplRecord 读取JPL记录
func readJplRecord(file *os.File, recordNum int) ([]Float64, error) {
	// 定位到记录位置
	offset := int64(recordNum * JplRecordSize)
	_, err := file.Seek(offset, io.SeekStart)
	if err != nil {
		return nil, err
	}
	
	// 读取记录
	buf := make([]byte, JplRecordSize)
	_, err = io.ReadFull(file, buf)
	if err != nil {
		return nil, err
	}
	
	// 转换为Float64数组
	reader := NewByteReader(buf, binary.LittleEndian)
	record := make([]Float64, JplNCoeff)
	
	for i := 0; i < JplNCoeff; i++ {
		record[i] = reader.ReadFloat64()
	}
	
	return record, nil
}

// computePosition 计算天体位置
func computePosition(record []Float64, et Float64, ntarg, ncent int, rrd *[6]Float64) error {
	var pv [15][6]Float64  // 增大数组以容纳所有JPL天体
	var t [2]Float64
	
	// 计算归一化时间
	s := et - record[0]
	t[0] = s / record[1]
	t[1] = 1.0
	
	// 计算各天体位置
	for k := 0; k < 11; k++ {
		// 获取切比雪夫多项式参数
		na := jplData.Header.IPT[0][k]
		nc := jplData.Header.IPT[1][k]
		ns := jplData.Header.IPT[2][k]
		
		if na == 0 {
			continue
		}
		
		// 计算子区间
		temp := Float64(ns) * t[0]
		l := int(temp)
		if l >= ns {
			l = ns - 1
		}
		
		tc := 2.0*(temp-Float64(l)) - 1.0
		
		// 计算切比雪夫多项式
		err := interpolateChebyshev(record, na+l*nc*3, nc, tc, pv[k][:])
		if err != nil {
			return err
		}
		
		// 缩放速度
		vfac := 2.0 * Float64(ns) / record[1]
		pv[k][3] *= vfac
		pv[k][4] *= vfac
		pv[k][5] *= vfac
	}
	
	// 处理特殊情况（如月球相对于地月系质心的位置）
	if ntarg == JMoon {
		// 月球位置 = 地月系质心位置 + 月球相对地月系质心位置
		for i := 0; i < 6; i++ {
			if JEmb < 15 && JMoon < 15 {
				pv[JMoon][i] = pv[JEmb][i] + pv[JMoon][i]
			}
		}
	}
	
	// 计算相对位置
	if ntarg == ncent {
		// 相对于自身
		for i := 0; i < 6; i++ {
			rrd[i] = 0.0
		}
	} else if ncent == JSbary {
		// 相对于太阳系质心
		for i := 0; i < 6; i++ {
			if ntarg < 15 {
				rrd[i] = pv[ntarg][i]
			}
		}
	} else {
		// 相对于其他天体
		for i := 0; i < 6; i++ {
			if ntarg < 15 && ncent < 15 {
				rrd[i] = pv[ntarg][i] - pv[ncent][i]
			}
		}
	}
	
	return nil
}

// interpolateChebyshev 切比雪夫多项式插值
func interpolateChebyshev(coeffs []Float64, offset, ncoeff int, x Float64, result []Float64) error {
	// 初始化
	for i := 0; i < 6; i++ {
		result[i] = 0.0
	}
	
	if ncoeff < 2 {
		return fmt.Errorf("系数数量太少: %d", ncoeff)
	}
	
	// 计算切比雪夫多项式
	t := make([]Float64, ncoeff)
	t[0] = 1.0
	t[1] = x
	
	for i := 2; i < ncoeff; i++ {
		t[i] = 2.0*x*t[i-1] - t[i-2]
	}
	
	// 计算位置（x, y, z）
	for j := 0; j < 3; j++ {
		for i := 0; i < ncoeff; i++ {
			result[j] += coeffs[offset+j*ncoeff+i] * t[i]
		}
	}
	
	// 计算速度（vx, vy, vz）
	if ncoeff >= 2 {
		dt := make([]Float64, ncoeff)
		dt[0] = 0.0
		dt[1] = 1.0
		
		for i := 2; i < ncoeff; i++ {
			dt[i] = 2.0*x*dt[i-1] + 2.0*t[i-1] - dt[i-2]
		}
		
		for j := 0; j < 3; j++ {
			for i := 0; i < ncoeff; i++ {
				result[j+3] += coeffs[offset+j*ncoeff+i] * dt[i]
			}
		}
	}
	
	return nil
}

// ByteReader 字节读取器
type ByteReader struct {
	data  []byte
	pos   int
	order binary.ByteOrder
}

// NewByteReader 创建字节读取器
func NewByteReader(data []byte, order binary.ByteOrder) *ByteReader {
	return &ByteReader{
		data:  data,
		pos:   0,
		order: order,
	}
}

// ReadFloat64 读取Float64
func (r *ByteReader) ReadFloat64() Float64 {
	if r.pos+8 > len(r.data) {
		return 0
	}
	
	val := r.order.Uint64(r.data[r.pos : r.pos+8])
	r.pos += 8
	return math.Float64frombits(val)
}

// ReadInt32 读取Int32
func (r *ByteReader) ReadInt32() Int32 {
	if r.pos+4 > len(r.data) {
		return 0
	}
	
	val := r.order.Uint32(r.data[r.pos : r.pos+4])
	r.pos += 4
	return Int32(val)
}

// IERSF5 坐标系转换（ICRS到FK5）
func IERSF5(xin [6]Float64, dir int) [6]Float64 {
	var xout [6]Float64
	
	// 简化的转换矩阵（实际应该使用精确的IERS矩阵）
	// 这里仅作示例，实际实现需要更精确的转换
	if dir > 0 {
		// ICRS到FK5
		for i := 0; i < 6; i++ {
			xout[i] = xin[i]
		}
	} else {
		// FK5到ICRS
		for i := 0; i < 6; i++ {
			xout[i] = xin[i]
		}
	}
	
	return xout
}

// 检查JPL文件是否可用
func IsJplAvailable() bool {
	return jplData.IsOpen
}

// 获取JPL文件信息
func GetJplInfo() (startJD, endJD, stepJD Float64, isOpen bool) {
	if jplData.IsOpen {
		return jplData.Header.StartJD, jplData.Header.EndJD, jplData.Header.StepJD, true
	}
	return 0, 0, 0, false
}