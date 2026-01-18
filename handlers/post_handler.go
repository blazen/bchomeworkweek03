package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"gin-examples/project/models"
	"gin-examples/project/services"
	"gin-examples/project/utils"
)

type PostHandler struct {
	postService *services.PostService
}

func NewPostHandler(postService *services.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

// 创建文章
func (h *PostHandler) CreatePost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req models.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.ParseValidationErrors(err))
		return
	}

	post, err := h.postService.CreatePost(userID.(uint), req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.Success(c, post)
}

// 查询用户的全部文章
func (h *PostHandler) ListPost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	pageNoStr := c.DefaultQuery("pageNo", "1") // 带默认值
	pageSizeStr := c.DefaultQuery("pageSize", "5")
	pageNo, err := strconv.Atoi(pageNoStr)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	posts, err := h.postService.ListPost(userID.(uint), pageNo, pageSize)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.Success(c, posts)
}

// 查询所有用户的文章
func (h *PostHandler) ListPostAll(c *gin.Context) {
	pageNoStr := c.DefaultQuery("pageNo", "1") // 带默认值
	pageSizeStr := c.DefaultQuery("pageSize", "5")
	pageNo, err := strconv.Atoi(pageNoStr)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	posts, err := h.postService.ListPostAll(pageNo, pageSize)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.Success(c, posts)
}

// 条件查询文章
func (h *PostHandler) ListPostByCondition(c *gin.Context) {
	var req models.ListPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.ParseValidationErrors(err))
		return
	}

	posts, err := h.postService.ListPostByCondition(req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.Success(c, posts)
}

// 主键查询文章
func (h *PostHandler) GetPostById(c *gin.Context) {
	id := c.Param("id")
	uintid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		fmt.Println("主键id字符串转 uint64 转换错误:", err)
		utils.HandleError(c, utils.NewAppError(409, "Invalid id"))
		return
	}

	post, err := h.postService.GetPostById(uint(uintid))
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.Success(c, post)
}

// 查询评论数量最多的文章
func (h *PostHandler) GetPostByMaxCommentNumber(c *gin.Context) {
	post, err := h.postService.GetPostByMaxCommentNumber()
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.Success(c, post)
}

// 更新用户的文章
func (h *PostHandler) UpdatePost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req models.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.ParseValidationErrors(err))
		return
	}

	post, err := h.postService.UpdatePost(userID.(uint), req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.Success(c, post)
}

// 删除用户的文章
func (h *PostHandler) DeletePost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id := c.Param("id")
	uintid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		fmt.Println("主键id字符串转 uint64 转换错误:", err)
		utils.HandleError(c, utils.NewAppError(409, "Invalid id"))
		return
	}

	// 64位系统不会导致溢出
	r, err := h.postService.DeletePost(userID.(uint), uint(uintid))
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.Success(c, r)
}
