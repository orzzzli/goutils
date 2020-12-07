package logger

import (
	"errors"
	"os"
	"time"
)

const (
	TimeStrFormat = "2006-01-02"
)

type Worker struct {
	idx uint
}

func newWorker(idx uint, workFunc func(worker *Worker)) *Worker {
	w := &Worker{
		idx: idx,
	}
	go workFunc(w)
	return w
}

type LogMessage struct {
	time    string
	message string
}

func newLogMessage(message string) *LogMessage {
	return &LogMessage{
		time:    time.Now().Format(TimeStrFormat),
		message: message,
	}
}

type Logger struct {
	path      string
	fileName  string
	separator byte
	fileFds   map[string]*os.File

	logChan   chan *LogMessage
	workerMap map[uint]*Worker

	managerError error
	workerError  error
}

/*
	bufferSize代表logChan缓存的message条数
	chunkSize 单位：KB
*/
func NewLogger(path string, fileName string, separator byte, bufferSize uint, workerSize uint) (*Logger, error) {
	l := &Logger{
		path:         path,
		fileName:     fileName,
		fileFds:      make(map[string]*os.File),
		separator:    separator,
		logChan:      make(chan *LogMessage, bufferSize),
		workerMap:    make(map[uint]*Worker),
		managerError: nil,
		workerError:  nil,
	}
	//打开两个文件避免文件切换时的锁
	//打开/创建当天文件
	nowTime := time.Now()
	nowPath := l.getFullPath(nowTime.Format(TimeStrFormat))
	f, err := openFile(nowPath)
	if err != nil {
		return nil, err
	}
	l.fileFds[nowTime.Format(TimeStrFormat)] = f
	//打开/创建下一个时间阶段文件
	nextTime := nowTime.Add(24 * time.Hour)
	nextPath := l.getFullPath(nextTime.Format(TimeStrFormat))
	f1, err := openFile(nextPath)
	if err != nil {
		return nil, err
	}
	l.fileFds[nextTime.Format(TimeStrFormat)] = f1

	go l.manage()
	var i uint
	for i = 0; i < workerSize; i++ {
		w := newWorker(i, l.worker)
		l.workerMap[i] = w
	}
	return l, nil
}

//写日志
func (l *Logger) Log(messages ...string) {
	logStr := ""
	for k, v := range messages {
		logStr += v
		if k != len(messages)-1 {
			logStr += string(l.separator)
		}
	}
	logStr += "\n"
	tempMessage := newLogMessage(logStr)
	l.logChan <- tempMessage
}

func (l *Logger) GetLastError() (error, error) {
	return l.managerError, l.workerError
}

//worker方法
func (l *Logger) worker(worker *Worker) {
	for {
		messageObj := <-l.logChan
		file, ok := l.fileFds[messageObj.time]
		if !ok {
			l.workerError = errors.New("cant found " + messageObj.time + " file")
			continue
		}
		_, err := file.Write([]byte(messageObj.message))
		if err != nil {
			l.workerError = err
		}
	}
}

//worker管理器
func (l *Logger) manage() {
	for {
		nowTime := time.Now()
		nextTime := nowTime.Add(24 * time.Hour)
		preTime := nowTime.Add(-(24 * time.Hour))
		//创建下一个文件
		_, ok := l.fileFds[nextTime.Format(TimeStrFormat)]
		if !ok {
			path := l.getFullPath(nextTime.Format(TimeStrFormat))
			f1, err := openFile(path)
			if err != nil {
				l.managerError = err
				goto CHECK_LAST
			}
			l.fileFds[nextTime.Format(TimeStrFormat)] = f1
		}
	CHECK_LAST:
		var f *os.File
		//为防止释放上一个文件，有协程还未写完，延迟最少三个小时再释放历史文件fd
		if nowTime.Hour() < 3 {
			goto GOSLEEP
		}
		//检查上一个文件
		f, ok = l.fileFds[preTime.Format(TimeStrFormat)]
		if ok {
			//释放
			err := f.Close()
			if err != nil {
				l.managerError = err
			}
			delete(l.fileFds, preTime.Format(TimeStrFormat))
		}
	GOSLEEP:
		time.Sleep(24 * time.Hour)
	}
}

//工具方法：获取当前完整路径
func (l *Logger) getFullPath(time string) string {
	fullPath := ""
	fullPath = l.path + time + "-" + l.fileName + ".log"
	return fullPath
}

//工具方法：获取24小时后完成路径
func (l *Logger) getNextPath() string {
	fullPath := ""
	nextTime := time.Now().Add(24 * time.Hour).Format("2006-01-02")
	fullPath = l.path + nextTime + "-" + l.fileName + ".log"
	return fullPath
}

//工具方法：追加模式打开文件
func openFile(path string) (*os.File, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return f, nil
}
