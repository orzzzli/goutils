package configer

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultScanSec     = 60
	DefaultHotLoading  = true
	DefaultDebug       = false
	DefaultCommentChar = ';'
)

type IniConfiger struct {
	path        string
	hotLoading  bool
	scanSec     int
	debug       bool
	commentChar byte

	lastMD5   string
	sections  map[string]map[string]string
	configMap map[string]string
}

func NewiniConfiger(path string) (*IniConfiger, error) {
	if path == "" {
		return nil, errors.New("[path] is empty")
	}
	exist, err := fileExists(path)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New("config file:" + path + " is not exist")
	}

	return &IniConfiger{
		path:        path,
		hotLoading:  DefaultHotLoading,
		scanSec:     DefaultScanSec,
		debug:       DefaultDebug,
		commentChar: DefaultCommentChar,

		sections:  make(map[string]map[string]string, 10),
		configMap: make(map[string]string, 10),
	}, nil
}

func (i *IniConfiger) SetCommentChar(commentChar byte) {
	i.commentChar = commentChar
}

func (i *IniConfiger) SetScanSec(scanSec int) {
	i.scanSec = scanSec
}

func (i *IniConfiger) SwitchHotLoading() {
	i.hotLoading = !i.hotLoading
}

func (i *IniConfiger) SwitchDebug() {
	i.debug = !i.debug
}

/*
	Read ini config and instantiate configer map.
*/
func (i *IniConfiger) Invoke() error {
	file, err := os.Open(i.path)
	if err != nil {
		return err
	}
	defer file.Close()

	i.sections = make(map[string]map[string]string, 10)
	i.configMap = make(map[string]string, 10)

	reader := bufio.NewReader(file)
	lastSection := ""
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		lineStr := string(line)

		exist, first, start := i.processComment(lineStr)
		//whole line is comment
		if exist && first {
			continue
		}

		section, sectionExist := i.processSection(lineStr, start)
		//whole line is section
		if sectionExist {
			_, ok := i.sections[section]
			//section already exist
			if ok {
				continue
			}
			i.sections[section] = make(map[string]string)
			lastSection = section
		}

		key, value, ok := i.processLine(lineStr, start)
		if ok {
			i.configMap[key] = value
			//dont have section
			if lastSection == "" {
				continue
			}
			i.sections[lastSection][key] = value
		}
	}
	return nil
}

/*
	Hot Loading config while file md5 changed.
*/
func (i *IniConfiger) hotLoadingConfiger() error {
	for {
		if i.hotLoading {
			file, err := os.Open(i.path)
			if err != nil {
				return errors.New("open config file err : " + err.Error())
			}
			md5Obj := md5.New()
			_, err = io.Copy(md5Obj, file)
			if err != nil {
				return errors.New("io copy file error : " + err.Error())
			}
			md5Str := hex.EncodeToString(md5Obj.Sum(nil))
			//first time
			if i.lastMD5 == "" {
				i.lastMD5 = md5Str
			} else if i.lastMD5 != md5Str { //config file changed
				err = i.Invoke()
				if err != nil {
					return err
				}
			}
			file.Close()
		}

		time.Sleep(time.Duration(i.scanSec) * time.Second)
	}
}

/*
	Get function return value and is ok.
*/
func (i *IniConfiger) GetString(section string, key string) (string, bool) {
	value, ok := i.sections[section][key]
	return value, ok
}

/*
	Get function return int and is ok.
	If not exist return -1 and false.
	If convert err return 0 and false.
*/
func (i *IniConfiger) GetInt(section string, key string) (int, bool) {
	value, ok := i.sections[section][key]
	if ok {
		valueInt, err := strconv.Atoi(value)
		if err == nil {
			return valueInt, true
		}
		return 0, false
	}
	return -1, false
}

/*
	Check file is dir and check file exist.
*/
func fileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		if info.IsDir() {
			return false, errors.New("path:" + path + " is a folder")
		}
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, nil
}

func (i *IniConfiger) processLine(context string, commentStart int) (key string, value string, ok bool) {
	start := strings.Index(context, "=")
	if start != -1 && start != 0 {
		//comment exist
		if commentStart != -1 {
			if commentStart < start {
				return "", "", false
			}
		}
		key = string([]rune(context)[:start])
		key = strings.Replace(key, " ", "", -1)
		if key == "" {
			return "", "", false
		}
		value = string([]rune(context)[start+1:])
		value = strings.Replace(value, " ", "", -1)
		if value == "" {
			return key, "", true
		}
		//check comment in value
		commentStart2 := strings.Index(value, string(i.commentChar))
		if commentStart2 != -1 && commentStart2 != 0 { //comment in middle of value
			value = string([]rune(value)[:commentStart2])
			return key, value, true
		} else if commentStart2 == 0 { //comment at first of value
			return key, "", true
		} else { //comment not exist
			return key, value, true
		}
	} else {
		return "", "", false
	}
}

func (i *IniConfiger) processSection(context string, commentStart int) (str string, exist bool) {
	section := ""
	start := strings.Index(context, "[")
	end := strings.Index(context, "]")
	if start != -1 && end != -1 {
		//comment exist
		if commentStart != -1 {
			if commentStart < start {
				return section, false
			}
		}
		section = string([]rune(context)[start+1 : end])
		return section, true
	} else {
		return section, false
	}
}

func (i *IniConfiger) processComment(context string) (exist bool, isFirst bool, index int) {
	start := strings.Index(context, string(i.commentChar))
	if start == 0 {
		return true, true, start
	} else if start == -1 {
		return false, false, start
	} else {
		return true, false, start
	}
}
