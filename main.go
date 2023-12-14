package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"ipfs-storage/dao"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/ipfs/go-ipfs-api"
)

var db *gorm.DB

func init() {
	db = dao.GetDB()
}
func main() {
	defer db.Close()

	db.AutoMigrate(&dao.FileHash{})
	//db.Model(&dao.FileHash{})
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	//r.GET("/", renderForm)
	r.GET("/upload", FileUploadFront)
	r.POST("/upload", handleFileUpload)
	r.GET("/GetFile", showUploadedFiles)
	//r.POST("/GetFile", showUploadedFiles)
	r.GET("/download", Handledownload)

	r.POST("/download", Handledownload)
	r.Run(":8081")
}
func HandledownloadFront(c *gin.Context) {
	c.HTML(200, "download.html", nil)
}
func Handledownload(c *gin.Context) {
	FileName := c.Query("FileName")
	ipfsHash := c.Query("hash")
	if ipfsHash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IPFS哈希不能为空"})
		return
	}
	shell := shell.NewShell("localhost:5001")
	downloadDir := "./downloads"

	// Download the file from IPFS
	filePath := filepath.Join(downloadDir, FileName)
	err := shell.Get(ipfsHash, filePath)
	log.Println(err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "从IPFS下载文件失败"})
		return
	}

	// Provide the downloaded file for download
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", FileName))
	c.File(filePath)
}
func showUploadedFiles(c *gin.Context) {
	// 从数据库中获取所有文件信息
	files, err := getAllFiles()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "handleupload.html", gin.H{"Error": "无法获取文件信息"})
		return
	}

	// 创建一个切片用于存储哈希与名字的对应关系
	var fileDetails []map[string]string
	for _, file := range files {
		fileDetail := map[string]string{
			"FileName": file.FileName,
			"IPFSHash": file.IPFSHash,
		}
		fileDetails = append(fileDetails, fileDetail)
	}

	// 渲染 handleupload.html 页面
	c.HTML(http.StatusOK, "handleupload.html", gin.H{
		"Error":       "",          // 如果有错误，可以在这里设置错误信息
		"FileDetails": fileDetails, // 传递哈希与名字的对应关系到模板
	})
}

func FileUploadFront(c *gin.Context) {
	c.HTML(http.StatusOK, "upload.html", nil)

}
func handleFileUpload(c *gin.Context) {
	c.HTML(http.StatusOK, "upload.html", nil)
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法获取文件"})
		return
	}

	// 保存文件到本地
	fileName := filepath.Join("uploads", file.Filename)
	err = c.SaveUploadedFile(file, fileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件保存失败"})
		return
	}

	// 将文件上传到IPFS
	hash, err := uploadToIPFS(fileName)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件上传到IPFS失败"})
		return
	}
	err = saveHashToDatabase(file.Filename, hash)
	// 返回上传成功的消息和IPFS哈希
	c.JSON(http.StatusOK, gin.H{"message": "文件上传成功", "ipfsHash": hash})
}

func uploadToIPFS(filePath string) (string, error) {
	// 创建IPFS API客户端
	shell := shell.NewShell("localhost:5001")

	// 通过IPFS API将文件添加到IPFS网络中
	hash, err := shell.AddDir(filePath)
	if err != nil {
		return "", err
	}

	fmt.Printf("文件成功上传到IPFS，IPFS哈希：%s\n", hash)
	return hash, nil
}

func saveHashToDatabase(fileName, ipfsHash string) error {
	// 连接数据库

	// 创建记录
	fileHash := dao.FileHash{

		FileName: fileName,
		IPFSHash: ipfsHash,
	}

	// 将记录保存到数据库
	if err := db.Create(&fileHash).Error; err != nil {
		return err
	}

	fmt.Println("IPFS哈希成功保存到数据库")
	return nil
}
func getAllFiles() ([]dao.FileHash, error) {
	var files []dao.FileHash

	// 查询数据库中的所有文件信息
	if err := db.Find(&files).Error; err != nil {
		return nil, err
	}

	return files, nil
}
