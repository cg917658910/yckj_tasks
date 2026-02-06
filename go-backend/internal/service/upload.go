package service

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"task-system-go/internal/config"

	"github.com/google/uuid"
)

func SaveImage(cfg config.Config, file *multipart.FileHeader) (map[string]string, error) {
	if file == nil {
		return nil, errors.New("请选择文件")
	}
	if file.Size > 5*1024*1024 {
		return nil, errors.New("文件过大")
	}
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".pdf": true, ".doc": true, ".docx": true}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowed[ext] {
		return nil, errors.New("不支持的文件类型")
	}
	name := uuid.New().String() + ext
	if err := os.MkdirAll(cfg.UploadDir, 0755); err != nil {
		return nil, errors.New("上传失败")
	}
	path := filepath.Join(cfg.UploadDir, name)
	if err := saveUploadedFile(file, path); err != nil {
		return nil, errors.New("上传失败")
	}
	// cfg.UploadDir 处理相对目录

	rel := filepath.ToSlash(filepath.Join(cfg.UploadDir, name))
	return map[string]string{
		"url":  cfg.BaseURL + "/" + rel,
		"path": rel,
	}, nil
}

// wrapped for easier testing
var saveUploadedFile = func(file *multipart.FileHeader, path string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}
