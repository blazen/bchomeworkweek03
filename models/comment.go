package models

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      // Foreign key to user
	PostID    uint      // Foreign key to user
	Content   string    `json:"content" gorm:"not null;size:100"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// CreatedAt utils.Time1    `json:"created_at"`		// 不能如此，创建报错
	// UpdatedAt utils.Time1    `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	User      User           `json:"-"`
}

type CreateCommentRequest struct {
	PostID  uint   `json:"post_id"`
	Content string `json:"content"`
}

// 评论创建钩子：文章评论数+1
func (a *Comment) AfterCreate(tx *gorm.DB) error {
	fmt.Printf("评论创建钩子执行 AfterCreate\n")
	postId := ContextValue(tx)
	fmt.Printf("评论创建钩子执行 AfterCreate 更新文章评论数+1 postId:%d", postId)
	var commentPost Post
	if err := tx.First(&commentPost, postId).Error; err != nil {
		return err
	}
	commentPost.CommentNumber++
	if 0 < commentPost.CommentNumber {
		commentPost.CommentStatus = "热评中"
	}
	if err := tx.Select("CommentNumber").Save(commentPost).Error; err != nil {
		return err
	}
	return nil
}

// func (a *Comment) AfterDelete(tx *gorm.DB) error {
// 	fmt.Printf("评论删除钩子执行 AfterDelete\n")
// 	postId := ContextValue(tx)
// 	fmt.Printf("评论删除钩子执行 AfterDelete 更新文章评论数-1 postId:%d", postId)
// 	var commentPost Post
// 	if err := tx.First(&commentPost, postId).Error; err != nil {
// 		return err
// 	}
// 	// 虽然类型为 uint 但实际数据库存储的可以是负数，所以要明确知道删除成功。
// 	commentPost.CommentNumber--
// 	if 0 == commentPost.CommentNumber {
// 		commentPost.CommentStatus = "无评论"
// 	}
// 	if err := tx.Save(commentPost).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

// 评论删除钩子：文章评论数-1
// 执行顺序： Delete > AfterDelete > next line of Delete function
func (a *Comment) AfterDelete(tx *gorm.DB) error {
	fmt.Printf("评论删除钩子执行 AfterDelete time=%v\n", time.Now())
	ctx := tx.Statement.Context
	if ctx == nil {
		ctx = context.Background() // 兜底：避免Context为nil
	}
	// time.Sleep(2 * time.Second) // obstruct next line of Delete function
	ra := tx.RowsAffected // 无法获取删除结果
	rowsAffected, ok := ctx.Value(CtxKeyRowsAffected).(*RowsAffectedContainer)
	fmt.Printf("评论删除钩子执行 AfterDelete 更新文章评论数-1 rowsAffected:%d, ra=%d, time=%v\n", rowsAffected, ra, time.Now())
	if !ok || rowsAffected == nil {
		fmt.Println("AfterDelete: 未获取到受影响行数")
		return nil
	}
	// 判断是否真的删除了数据
	if rowsAffected.Value > 0 {
		fmt.Printf("AfterDelete: 数据删除成功，受影响行数=%d\n", rowsAffected.Value)

		// 计算文章评论数
		postId := rowsAffected.PostId
		var commentPost Post
		if err := tx.First(&commentPost, postId).Error; err != nil {
			return err
		}
		// 虽然类型为 uint 但实际数据库存储的可以是负数，所以要明确知道删除成功。
		if 0 < commentPost.CommentNumber {
			commentPost.CommentNumber--
		} else {
			commentPost.CommentNumber = 0
			commentPost.CommentStatus = "无评论"
		}
		if err := tx.Select("CommentNumber").Save(commentPost).Error; err != nil {
			return err
		}
	} else {
		fmt.Printf("AfterDelete: 未删除任何内容\n")
	}
	return nil
}
