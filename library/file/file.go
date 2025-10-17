package file

import (
	"archive/zip"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	httpClient "your_project/library/http"
)

// ========================================
// 文件下载
// ========================================

// Download 下载文件到指定路径
func Download(url, savePath string) error {
	resp, err := httpClient.Client.R().
		SetOutput(savePath).
		Get(url)

	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("download failed: status code %d", resp.StatusCode())
	}

	return nil
}

// DownloadWithProgress 下载文件并显示进度
func DownloadWithProgress(url, savePath string, callback func(current, total int64)) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("download failed: status code %d", resp.StatusCode)
	}

	out, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer out.Close()

	total := resp.ContentLength
	var current int64 = 0

	buffer := make([]byte, 32*1024) // 32KB buffer
	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			out.Write(buffer[:n])
			current += int64(n)
			if callback != nil {
				callback(current, total)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	return nil
}

// DownloadToBytes 下载文件到字节数组
func DownloadToBytes(url string) ([]byte, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	if !httpClient.IsSuccess(resp) {
		return nil, fmt.Errorf("download failed: status code %d", resp.StatusCode())
	}

	return resp.Body(), nil
}

// ========================================
// 文件操作
// ========================================

// Exists 判断文件或目录是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// IsFile 判断是否为文件
func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// IsDir 判断是否为目录
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// FileSize 获取文件大小（字节）
func FileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// FileSizeFormat 格式化文件大小
func FileSizeFormat(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// CreateDir 创建目录（支持多级）
func CreateDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// Copy 复制文件
func Copy(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// Move 移动文件
func Move(src, dst string) error {
	return os.Rename(src, dst)
}

// Delete 删除文件或目录
func Delete(path string) error {
	return os.RemoveAll(path)
}

// ReadFile 读取文件内容
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// WriteFile 写入文件
func WriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

// AppendFile 追加内容到文件
func AppendFile(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	return err
}

// ========================================
// 文件类型检测
// ========================================

// GetMimeType 获取文件 MIME 类型
func GetMimeType(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	return http.DetectContentType(buffer[:n]), nil
}

// GetExtension 获取文件扩展名（带点）
func GetExtension(path string) string {
	return filepath.Ext(path)
}

// GetExtensionWithoutDot 获取文件扩展名（不带点）
func GetExtensionWithoutDot(path string) string {
	ext := filepath.Ext(path)
	if len(ext) > 0 {
		return ext[1:]
	}
	return ""
}

// GetFileName 获取文件名（不含扩展名）
func GetFileName(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext)
}

// GetFileNameWithExt 获取文件名（含扩展名）
func GetFileNameWithExt(path string) string {
	return filepath.Base(path)
}

// GetDir 获取文件所在目录
func GetDir(path string) string {
	return filepath.Dir(path)
}

// ========================================
// 文件压缩与解压
// ========================================

// ZipFiles 压缩文件
func ZipFiles(zipPath string, files []string) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, file := range files {
		if err := addFileToZip(zipWriter, file); err != nil {
			return err
		}
	}

	return nil
}

// addFileToZip 添加文件到 zip
func addFileToZip(zipWriter *zip.Writer, filename string) error {
	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = filepath.Base(filename)
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, fileToZip)
	return err
}

// UnzipFile 解压文件
func UnzipFile(zipPath, destDir string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		path := filepath.Join(destDir, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		// 创建文件目录
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		// 创建文件
		destFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		fileInArchive, err := file.Open()
		if err != nil {
			destFile.Close()
			return err
		}

		if _, err := io.Copy(destFile, fileInArchive); err != nil {
			destFile.Close()
			fileInArchive.Close()
			return err
		}

		destFile.Close()
		fileInArchive.Close()
	}

	return nil
}

// ========================================
// 文件哈希
// ========================================

// MD5File 计算文件 MD5
func MD5File(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// ========================================
// 批量操作
// ========================================

// ListFiles 列出目录下所有文件
func ListFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

// ListFilesByExt 列出指定扩展名的文件
func ListFilesByExt(dir, ext string) ([]string, error) {
	var files []string

	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ext {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

// DeleteFilesByExt 删除指定扩展名的文件
func DeleteFilesByExt(dir, ext string) error {
	files, err := ListFilesByExt(dir, ext)
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := os.Remove(file); err != nil {
			return err
		}
	}

	return nil
}

// ========================================
// 使用示例（注释）
// ========================================
/*
// 1. 下载文件
err := file.Download("https://example.com/file.pdf", "/path/to/save/file.pdf")

// 2. 带进度的下载
err := file.DownloadWithProgress(url, savePath, func(current, total int64) {
    progress := float64(current) / float64(total) * 100
    fmt.Printf("下载进度: %.2f%%\n", progress)
})

// 3. 复制文件
err := file.Copy("/path/to/source.txt", "/path/to/dest.txt")

// 4. 压缩文件
files := []string{"file1.txt", "file2.txt", "file3.txt"}
err := file.ZipFiles("archive.zip", files)

// 5. 解压文件
err := file.UnzipFile("archive.zip", "/path/to/extract")

// 6. 计算文件 MD5
md5, err := file.MD5File("/path/to/file.txt")
fmt.Println("MD5:", md5)

// 7. 获取文件信息
size, _ := file.FileSize("/path/to/file.txt")
fmt.Println("Size:", file.FileSizeFormat(size))

mimeType, _ := file.GetMimeType("/path/to/file.pdf")
fmt.Println("MIME:", mimeType)

// 8. 列出目录下所有文件
files, err := file.ListFiles("/path/to/dir")
for _, f := range files {
    fmt.Println(f)
}

// 9. 批量删除指定类型文件
err := file.DeleteFilesByExt("/path/to/dir", ".tmp")
*/
