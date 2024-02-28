/*
* @Author: Oatmeal107
* @Date:   2023/6/16 16:45
 */

package utils

import (
	"github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"os"
	"strconv"
)

// UploadAvatarLocal 上传头像到本地
func UploadAvatarLocal(file multipart.File, fileName string, userId uint) (filePath string, err error) {
	idStr := strconv.Itoa(int(userId))
	basePath := "./resources/static/imgs/avatar/" + idStr + "/"
	err = os.RemoveAll(basePath)
	if err != nil {
		logrus.Infoln(err)
	}
	if !DirExistOrNot(basePath) {
		CreateDir(basePath)
	}
	avatarPath := basePath + fileName
	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	err = os.WriteFile(avatarPath, content, 0666)
	if err != nil {
		return "", err
	}
	return avatarPath, err
}

func DirExistOrNot(fileAddr string) bool { // todo 这个函数有点问题,在文件夹不存在的时候会出一个Info
	s, err := os.Stat(fileAddr)
	if err != nil {
		logrus.Infoln(err)
		return false
	}
	return s.IsDir()
}

func CreateDir(fileAddr string) bool {
	err := os.Mkdir(fileAddr, 0777)
	if err != nil {
		logrus.Infoln(err)
		return false
	}
	return true
}
