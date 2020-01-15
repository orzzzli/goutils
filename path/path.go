package path

import (
	"os"
	"strings"
)

/*
	传入路径获取文件名称。(最后一格内容)
*/
func GetFileName(path string) string {
	tempList := strings.Split(path,string(os.PathSeparator))
	return tempList[len(tempList) - 1]
}

/*
	获取两个路径的差值
*/
func GetPathDiff(base string, other string, withFile bool) string {
	diffPath := strings.Replace(base,other,"",-1)
	if withFile {
		return diffPath
	}else{
		tempList := strings.Split(diffPath,string(os.PathSeparator))
		tempStr := ""
		for k,v := range tempList {
			if k == len(tempList)-1 {
				break
			}
			tempStr += string(os.PathSeparator)+v
		}
		return tempStr[1:]
	}
}

/*
	获取目录下全部文件列表
*/
func GetAllFiles(path string, list *[]string, withDir bool) error {
	f, err := os.OpenFile(path, os.O_RDONLY, os.ModeDir)
	if err != nil {
		return err
	}
	defer f.Close()
	fileInfo, _ := f.Readdir(-1)
	separator := string(os.PathSeparator)

	for _, info := range fileInfo {
		if info.IsDir() {
			if withDir {
				*list = append(*list, path + separator + info.Name())
			}
			GetAllFiles(path + separator + info.Name(),list,withDir)
		} else {
			*list = append(*list, path + separator + info.Name())
		}
	}

	return nil
}