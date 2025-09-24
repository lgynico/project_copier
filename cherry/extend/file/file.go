package cherryFile

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	cerr "github.com/lgynico/project-copier/cherry/error"
	cslice "github.com/lgynico/project-copier/cherry/extend/slice"
)

func JudgeFile(filePath string) (string, bool) {
	if filePath == "" {
		return filePath, false
	}

	var p, n string
	index := strings.LastIndex(filePath, "/")
	if index > 0 {
		p = filePath[:index]
		n = filePath[index+1:]
	} else {
		p = "./"
		n = filePath
	}

	newPath, found := JudgePath(p)
	if !found {
		return "", false
	}

	fullFilePath := filepath.Join(newPath, n)
	if IsFile(fullFilePath) {
		return fullFilePath, true
	}

	return "", false
}

func JudgePath(filePath string) (string, bool) {
	tmpPath := filepath.Join(GetWorkDir(), filePath)
	if IsDir(tmpPath) {
		return tmpPath, true
	}

	tmpPath = filepath.Join(GetCurrentDirectory(), filePath)
	if IsDir(tmpPath) {
		return tmpPath, true
	}

	dir := GetStackDir()
	for _, d := range dir {
		tmpPath = filepath.Join(d, filePath)
		if IsDir(tmpPath) {
			return tmpPath, true
		}
	}

	if IsDir(filePath) {
		return filePath, true
	}

	return "", false
}

func IsDir(dirPath string) bool {
	info, err := os.Stat(dirPath)
	return err == nil && info.IsDir()
}

func IsFile(fullPath string) bool {
	info, err := os.Stat(fullPath)
	return err == nil && !info.IsDir()
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	return strings.Replace(dir, "\\", "/", -1)
}

func GetCurrentPath() string {
	var absPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		absPath = path.Dir(filename)
	}

	return absPath
}

// TODO 结构和语法可以优化下
func GetStackDir() []string {
	var dir []string
	var buf [2 << 16]byte
	stack := string(buf[:runtime.Stack(buf[:], true)])
	lines := strings.Split(strings.TrimSpace(stack), "\n")

	for _, line := range lines {
		lastLine := strings.TrimSpace(line)
		lastIndex := strings.LastIndex(lastLine, "/")
		if lastIndex < 1 {
			continue
		}

		thisDir := lastLine[:lastIndex]
		if _, err := os.Stat(thisDir); err != nil {
			continue
		}

		if _, ok := cslice.StringIn(thisDir, dir); ok {
			continue
		}

		dir = append(dir, thisDir)
	}

	return dir
}

func GetWorkDir() string {
	p, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return p
}

func JoinPath(elem ...string) (string, error) {
	filePath := filepath.Join(elem...)
	if err := CheckPath(filePath); err != nil {
		return filePath, err
	}

	return filePath, nil
}

func CheckPath(filePath string) error {
	if _, err := os.Stat(filePath); err != nil {
		return err
	}

	return nil
}

func GetFileName(filePath string, removeExt bool) string {
	fileName := path.Base(filePath)
	if !removeExt {
		return fileName
	}

	suffix := path.Ext(fileName)
	return strings.TrimSuffix(fileName, suffix)
}

func WalkFiles(rootPath string, fileSuffix string) []string {
	var files []string

	rootPath, found := JudgePath(rootPath)
	if !found {
		return files
	}

	filepath.Walk(rootPath, func(path string, info fs.FileInfo, err error) error {
		if fileSuffix != "" && !strings.HasSuffix(path, fileSuffix) {
			return nil
		}

		files = append(files, path)
		return nil
	})

	return files
}

func ReadDir(rootPath string, filePrefix, fileSuffix string) ([]string, error) {
	var files []string

	rootPath, found := JudgePath(rootPath)
	if !found {
		return files, cerr.Errorf("path = %s, file not found.", rootPath)
	}

	fileInfo, err := os.ReadDir(rootPath)
	if err != nil {
		return nil, err
	}

	for _, info := range fileInfo {
		if info.IsDir() {
			continue
		}

		if filePrefix != "" && !strings.HasPrefix(info.Name(), filePrefix) {
			continue
		}

		if fileSuffix != "" && !strings.HasSuffix(info.Name(), fileSuffix) {
			continue
		}

		files = append(files, info.Name())
	}

	return files, nil
}
