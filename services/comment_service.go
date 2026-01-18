package services

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"gin-examples/project/models"
	"gin-examples/project/utils"
)

type CommentService struct {
	db *gorm.DB
}

func NewCommentService(db *gorm.DB) *CommentService {
	return &CommentService{db: db}
}

// 创建评论
func (s *CommentService) CreateComment(userId uint, req models.CreateCommentRequest) (*models.Comment, error) {
	// 检查文章是否已存在
	var existingPost models.Post
	if err := s.db.First(&existingPost, req.PostID).Error; err != nil {
		return nil, utils.NewAppError(409, "Post not exist")
	}

	ctx := models.ContextWithValue(req.PostID)

	comment := models.Comment{
		PostID:  req.PostID,
		UserID:  userId,
		Content: req.Content,
	}

	if err := s.db.WithContext(ctx).Create(&comment).Error; err != nil {
		return nil, err
	}

	return &comment, nil
}

// 查询文章的全部评论（查询文章时，通过 preload 可以自动关联查询出评论。评论分页需要继续使用此函数。）
func (s *CommentService) ListCommentByPostId(postId uint, pageNo, pageSize int) ([]models.Comment, error) {
	var comments []models.Comment
	if err := s.db.Scopes(utils.Sql.Paginate(pageNo, pageSize)).
		Where("post_id = ?", postId).
		Order("created_at desc").
		Find(&comments).Error; err != nil {
		return nil, utils.NewAppError(409, "Post Comment not exist")
	}

	return comments, nil
}

// 删除评论（使用 Unscoped 物理删除）
// func (s *CommentService) DeleteComment(userId uint, postId uint, id uint) (bool, error) {
// 	ctx := models.ContextWithValue(postId)
// 	// 如此传入错误的参数也会执行，且没法判断是否真的执行了删除。
// 	// 只有在数据库层面出错时才会非 nil（如网络断开、表不存在），ID 不存在不属于错误。
// 	// if err := s.db.WithContext(ctx).Where("id = ? and user_id = ? and post_id=?", id, userId, postId).Unscoped().Delete(&models.Comment{}).Error; err != nil {
// 	// 	return false, utils.NewAppError(409, "Comment delete failed")
// 	// }

// 	// 判断删除的行数 RowsAffected >0
// 	result := s.db.WithContext(ctx).Where("id = ? and user_id = ? and post_id=?", id, userId, postId).Unscoped().Delete(&models.Comment{})
// 	if result.Error != nil {
// 		return false, result.Error
// 	}
// 	if result.RowsAffected > 0 {
// 		return true, nil // 真的删除了数据
// 	} else {
// 		return false, nil // ID不存在，无数据被删除
// 	}
// }

func (s *CommentService) DeleteComment(userId uint, postId uint, id uint) (bool, error) {
	// 查询要删除的数据
	fmt.Printf("删除 userId=%d, postId=%d, id=%d \n", userId, postId, id)
	var count int64
	if err := s.db.Model(&models.Comment{}).Where("id = ? and user_id = ? and post_id=?", id, userId, postId).Count(&count).Error; err != nil {
		fmt.Println("删除错误")
		return false, err
	}
	fmt.Printf("删除 userId=%d, postId=%d, id=%d, count=%d \n", userId, postId, id, count)

	// context
	container := &models.RowsAffectedContainer{
		Value:  count,
		PostId: postId,
	}
	originalCtx := s.db.Statement.Context
	if originalCtx == nil {
		originalCtx = context.Background()
	}
	ctx := context.WithValue(originalCtx, models.CtxKeyRowsAffected, container)
	s.db.Statement.Context = ctx

	// physical delete by Unscoped
	fmt.Printf("DeleteComment time=%v\n", time.Now())
	result := s.db.Where("id = ? and user_id = ? and post_id=?", id, userId, postId).Unscoped().Delete(&models.Comment{})
	if result.Error != nil {
		return false, result.Error
	}

	// 判断删除的行数 RowsAffected >0
	container.Value = result.RowsAffected
	fmt.Printf("result.RowsAffected=%d,time=%v\n", result.RowsAffected, time.Now())

	if result.RowsAffected > 0 {
		return true, nil // 真的删除了数据
	} else {
		return false, nil // ID不存在，无数据被删除
	}
}
