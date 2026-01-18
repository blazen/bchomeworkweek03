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

type CommentHandler struct {
	commentService *services.CommentService
}

func NewCommentHandler(commentService *services.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

// 创建评论
func (h *CommentHandler) CreateComment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.ParseValidationErrors(err))
		return
	}

	comment, err := h.commentService.CreateComment(userID.(uint), req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.Success(c, comment)
}

// 查询文章的全部评论（查询文章时，通过 preload 可以自动关联查询出评论。评论分页需要继续使用此函数。）
func (h *CommentHandler) ListCommentByPostId(c *gin.Context) {
	id := c.Param("postId")
	uintid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		fmt.Println("主键id字符串转 uint64 转换错误:", err)
		utils.HandleError(c, utils.NewAppError(409, "Invalid id"))
		return
	}
	pageNo, pageSize := utils.GetQueryPage(c)
	comments, err := h.commentService.ListCommentByPostId(uint(uintid), pageNo, pageSize)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.Success(c, comments)
}

// 删除评论
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	postId := c.Param("postId")
	uintpostId, err := strconv.ParseUint(postId, 10, 64)
	if err != nil {
		fmt.Println("主键id字符串转 uint64 转换错误:", err)
		utils.HandleError(c, utils.NewAppError(409, "Invalid id"))
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
	r, err := h.commentService.DeleteComment(userID.(uint), uint(uintpostId), uint(uintid))
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.Success(c, r)
}
