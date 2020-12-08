package snowflake

import (
	"container/list"
	"errors"
	"strconv"
	"sync"
	"time"
)

type RingBuffer struct {
	li     *list.List
	active *list.Element
}

func newRing() *RingBuffer {
	return &RingBuffer{
		li:     list.New(),
		active: nil,
	}
}

func (r *RingBuffer) pushBack(v int64) {
	r.li.PushBack(v)
}

func (r *RingBuffer) resetActive() {
	e := r.li.Front()
	r.active = e
}

//获取active的值并指针后移
func (r *RingBuffer) moveToNext() (int64, bool) {
	if r.active == nil {
		return 0, true
	}
	res := r.active.Value
	e := r.active.Next()
	r.active = e
	return res.(int64), false
}

/*
	高位第1位，永远为0，标识符号
	41位时间差，约等于69年
	1位数据中心id，0-1
	9位机器id，0-511
	12位每毫秒序列，0-4095
*/
type SnowFlake struct {
	mask           int64
	startTimeStamp int64
	dataCenterBit  uint64
	workerIdBit    uint64
	seqBit         uint64

	dataCenterId int
	workerId     int
	sequence     int64
	maxSeq       int64

	lastTimeStamp int64
	mu            sync.Mutex

	ring *RingBuffer
}

func NewSnowFlake(startTimeStamp int64, dataCenterBit uint64, workerIdBit uint64, seqBit uint64, dataCenterId int, workerId int) (*SnowFlake, error) {
	if dataCenterBit+workerIdBit+seqBit+42 > 64 {
		return nil, errors.New("total bit larger than 64 bit")
	}
	maxDataCenterId := 1 << dataCenterBit
	if dataCenterId > maxDataCenterId {
		return nil, errors.New("dataCenterId larger than " + strconv.Itoa(maxDataCenterId))
	}
	maxWorkId := 1 << workerIdBit
	if workerId > maxWorkId {
		return nil, errors.New("workerId larger than " + strconv.Itoa(maxWorkId))
	}
	maxSeq := 1 << seqBit
	r := newRing()
	for i := 0; i < maxSeq; i++ {
		r.pushBack(int64(i))
	}
	r.resetActive()
	return &SnowFlake{
		mask:           0xFFFFFFFF,
		startTimeStamp: startTimeStamp,
		dataCenterBit:  dataCenterBit,
		workerIdBit:    workerIdBit,
		seqBit:         seqBit,

		dataCenterId: dataCenterId,
		workerId:     workerId,
		sequence:     0,
		maxSeq:       int64(maxSeq),

		lastTimeStamp: 0,

		ring: r,
	}, nil
}

//加锁版本
func (s *SnowFlake) Get() int64 {
RETRY:
	s.mu.Lock()
	defer s.mu.Unlock()
	curTimeStamp := time.Now().UnixNano() / 1e6
	//需要顺序执行
	if curTimeStamp == s.lastTimeStamp {
		s.sequence++
		if s.sequence > s.maxSeq {
			time.Sleep(1 * time.Millisecond)
			s.mu.Unlock()
			goto RETRY
		}
	} else {
		s.lastTimeStamp = curTimeStamp
		s.sequence = 0
	}
	var timeGap, dataCenter, worker int64
	timeGap = (curTimeStamp - s.startTimeStamp) << (s.dataCenterBit + s.workerIdBit + s.seqBit)
	dataCenter = int64(s.dataCenterId) << (s.workerIdBit + s.seqBit)
	worker = int64(s.workerId) << s.seqBit
	return timeGap | dataCenter | worker | s.sequence
}

func (s *SnowFlake) GetFromRing() int64 {
RETRY:
	s.mu.Lock()
	curTimeStamp := time.Now().UnixNano() / 1e6
	var seq int64
	//需要顺序执行
	if curTimeStamp == s.lastTimeStamp {
		seq1, end := s.ring.moveToNext()
		if end {
			time.Sleep(1 * time.Millisecond)
			s.mu.Unlock()
			goto RETRY
		}
		seq = seq1
	} else {
		s.lastTimeStamp = curTimeStamp
		s.ring.resetActive()
		seq1, _ := s.ring.moveToNext()
		seq = seq1
	}
	s.mu.Unlock()
	var timeGap, dataCenter, worker int64
	timeGap = (curTimeStamp - s.startTimeStamp) << (s.dataCenterBit + s.workerIdBit + s.seqBit)
	dataCenter = int64(s.dataCenterId) << (s.workerIdBit + s.seqBit)
	worker = int64(s.workerId) << s.seqBit
	return timeGap | dataCenter | worker | seq
}
