/*
-------------------------------------------------
   Author :       zlyuan
   date：         2019/5/14
   Description :
-------------------------------------------------
*/

package zpbar

import (
    "fmt"
    "strings"
    "sync/atomic"
    "time"
)

const (
    pbarWidth      = 50     //进度条字符宽度
    lowerSpeed     = 5    //低速度阈值
    updateInterval = 0.01e9 //速度计算间隔时间
    drawInterval   = 0.3e9  //渲染间隔
    oneSecond      = 1e9    //一秒
)

type PBar struct {
    total int64  //目标值
    desc  string //描述
    unit  string //单位
    color ColorType

    count   int64 //计数
    run     int32 //运行中
    visible int32 //是否显示

    countInterval int64 //最后两次计数平均时间间隔
    timerCount    int64 //每秒计数
    startTime     int64 //开始时间
}

//创建一个进度条
func NewPbar(options ...Option) *PBar {
    b := &PBar{
        desc:    "zpbar",
        unit:    "it",
        visible: 1,
    }
    for _, o := range options {
        o(b)
    }
    return b
}

type Option func(*PBar)

//开始进度条
func (m *PBar) Start() {
    if !atomic.CompareAndSwapInt32(&m.run, 0, 1) {
        return
    }

    m.countInterval = oneSecond
    m.timerCount = 0
    m.startTime = time.Now().UnixNano()

    go func() {
        var oldCount = m.count
        for atomic.LoadInt32(&m.run) == 1 {
            time.Sleep(oneSecond)

            nowCount := atomic.LoadInt64(&m.count)
            atomic.StoreInt64(&m.timerCount, nowCount-oldCount)
            oldCount = nowCount
        }
    }()

    go func() {
        m.drawStart()
        for atomic.LoadInt32(&m.run) == 1 {
            m.draw()
            time.Sleep(drawInterval)
        }
    }()

    go func() {
        var oldCount = atomic.LoadInt64(&m.count)
        var oldCountTime = m.startTime

        for atomic.LoadInt32(&m.run) == 1 {
            time.Sleep(updateInterval)
            nowCount := atomic.LoadInt64(&m.count)
            if nowCount != oldCount {
                now := time.Now().UnixNano()
                atomic.StoreInt64(&m.countInterval, (now-oldCountTime)/(nowCount-oldCount)) //计算时间间隔
                oldCount = nowCount
                oldCountTime = now
            }
        }
    }()
}

//关闭进度条
func (m *PBar) Close() {
    m.draw()
    atomic.StoreInt32(&m.run, 0)
    m.drawEnd()
}

//是否运行中
func (m *PBar) IsRun() bool {
    return atomic.LoadInt32(&m.run) == 1
}

//添加完成量
func (m *PBar) Add(c int64) {
    atomic.AddInt64(&m.count, c)
}

//完成量+1
func (m *PBar) Done() {
    atomic.AddInt64(&m.count, 1)
}

//获取已完成量
func (m *PBar) Count() int64 {
    return atomic.LoadInt64(&m.count)
}

//获取目标量
func (m *PBar) Total() int64 {
    return m.total
}

//设置进度条是否显示
func (m *PBar) SetVisible(show bool) {
    if show {
        if atomic.CompareAndSwapInt32(&m.visible, 0, 1) {
            if m.IsRun() {
                printOs(getColorStartFlag(m.color))
            }
        }
    } else {
        if atomic.CompareAndSwapInt32(&m.visible, 1, 0) {
            if m.IsRun() {
                printOs(getColorEndFlag())
                printOs("\n")
            }
        }
    }
}

//获取进度条是否显示
func (m *PBar) Visible() bool {
    return atomic.LoadInt32(&m.visible) == 1
}

//获取进度条描述
func (m *PBar) Desc() string {
    return m.desc
}

//获取单位
func (m *PBar) Unit() string {
    return m.unit
}

//获取进度条颜色
func (m *PBar) Color() ColorType {
    return m.color
}

func (m *PBar) drawStart() {
    if m.Visible() {
        printOs(getColorStartFlag(m.color))
    }
}

func (m *PBar) drawEnd() {
    if m.Visible() {
        printOs(getColorEndFlag())
        printOs("\n")
    }
}

func (m *PBar) draw() {
    if m.Visible() && m.IsRun() {
        var text = m.toString()
        printOs("\r")
        printOs(text)
    }
}

func (m *PBar) toString() string {
    var allCount = atomic.LoadInt64(&m.total)
    var nowCount = atomic.LoadInt64(&m.count)
    var countInterval = atomic.LoadInt64(&m.countInterval)
    var timerCount = atomic.LoadInt64(&m.timerCount)

    var now = time.Now().UnixNano()
    var runTime = timeFmt((now - m.startTime) / oneSecond) //运行时间
    var needCount = allCount - nowCount                    //剩余数量

    var speed string
    var needTime string

    if countInterval >= oneSecond {
        var s = float64(countInterval) / float64(oneSecond)
        speed = fmt.Sprintf("%.2fs/%s", s, m.unit) //秒每条
        needTime = timeFmt(int64(float64(needCount) * s))
    } else if timerCount > lowerSpeed {
        speed = fmt.Sprintf("%d%s/s", timerCount, m.unit) //条每秒
        needTime = timeFmt(needCount / timerCount)
    } else {
        var c = float64(oneSecond) / float64(countInterval)
        speed = fmt.Sprintf("%.2f%s/s", c, m.unit) //条每秒
        needTime = timeFmt(int64(float64(needCount) / c))
    }

    if needCount < 0 || allCount == 0 {
        //描述  数值单位 [总时间, 速度]
        return fmt.Sprintf("%s %2d%s [%s, %s]", m.desc, nowCount, m.unit, runTime, speed)
    }

    var rate = float64(nowCount) / float64(allCount)

    var bar string
    var barlen = int(rate * pbarWidth)
    bar = fmt.Sprintf("%s%s", strings.Repeat(">", barlen), strings.Repeat(" ", pbarWidth-barlen))

    //描述 百分比% [进度条] 数值/总量 [总时间<剩余时间, 速度]
    return fmt.Sprintf("%s %3d%% [%s] %2d/%d [%s<%s, %s]", m.desc, int(rate*100), bar, nowCount, allCount, runTime, needTime, speed)

}

// 设置最大值
func WithTotal(total int64) Option {
    return func(pbar *PBar) {
        pbar.total = total
    }
}

// 设置描述
func WithDesc(desc string) Option {
    return func(pbar *PBar) {
        pbar.desc = desc
    }
}

// 设置数值单位
func WithUnit(unit string) Option {
    return func(pbar *PBar) {
        pbar.unit = unit
    }
}

// 设置进度条颜色
func WithColor(color ColorType) Option {
    return func(pbar *PBar) {
        pbar.color = color
    }
}

// 设置开始时数值
func WithStartCount(c int64) Option {
    return func(pbar *PBar) {
        pbar.count = c
    }
}
