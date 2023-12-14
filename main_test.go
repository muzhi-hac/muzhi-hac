package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUploadToIPFS(t *testing.T) {
	// 测试上传到IPFS的函数
	filePath := "/home/wang/GolandProjects/ipfs-storage/uploads/2_Space Complexity.pptx"
	hash, err := uploadToIPFS(filePath)

	assert.NoError(t, err, "Expected no error")
	assert.NotEmpty(t, hash, "Expected non-empty hash")
}

func TestSaveHashToDatabase(t *testing.T) {
	// 测试保存到数据库的函数
	fileName := "/home/wang/GolandProjects/ipfs-storage/uploads/2_Space Complexity.pptx"
	ipfsHash := "QmQnhUn3hYCE7mQ8txdjGgWR8Sc99QxrhWz929SB9U59zC"

	err := saveHashToDatabase(fileName, ipfsHash)

	assert.NoError(t, err, "Expected no error")
	// TODO: Add additional assertions or queries to check if the data is saved correctly in the database
}
func TestHandledownload(t *testing.T) {
	// 创建一个虚拟的gin.Context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 设置下载请求的参数，例如IPFS哈希
	c.Request, _ = http.NewRequest(http.MethodGet, "/download?hash=QmQnhUn3hYCE7mQ8txdjGgWR8Sc99QxrhWz929SB9U59zC&FileName=Assignment3.md", nil)

	// 调用Handledownload函数
	Handledownload(c)

	// 检查HTTP响应状态码
	assert.Equal(t, http.StatusOK, w.Code)

	// 在实际测试中，你可能还需要进一步检查下载的文件内容或其他方面的期望结果。
	// 请根据你的具体需求进行扩展。
}

// 如果需要更复杂的测试，你可能需要模拟数据库或使用测试数据库。
