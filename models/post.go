package models

import (
	"context"
	"fmt"
	"time"

	"gin-examples/project/utils"

	"gorm.io/gorm"
)

// 输入审计
type auditInputFields struct {
	AuditBy      string // 审计者，机器人或者人工
	AuditVersion string // 审计版本
	AuditStatus  string // 审计状态
}

type Post struct {
	ID     uint `json:"id" gorm:"primaryKey"`
	UserID uint // Foreign key to user
	// Title         string           `json:"title" gorm:"not null;size:50;uniqueIndex"` // 追加 uniqueIndex 时，若是数据中存在重复的数据，那么启动失败。
	Title         string           `json:"title" gorm:"not null;size:50"` // 加上的索引，要想去掉，只能手动去掉（删除数据库文件，重新创建）
	Content       string           `json:"content" gorm:"not null;size:100"`
	CommentNumber uint             `json:"comment_number" gorm:"default:0"`
	CommentStatus string           `json:"comment_status"`
	CreatedAt     utils.Time1      `json:"created_at"`
	UpdatedAt     utils.Time1      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt   `json:"-" gorm:"index"`
	Comments      []Comment        `json:"comments"`
	User          User             `json:"-"`
	Audit         auditInputFields `json:"-" gorm:"embedded"`
}

type CreatePostRequest struct {
	Title   string `json:"title" gorm:"not null;size:50"`
	Content string `json:"content" gorm:"not null;size:100"`
}

type UpdatePostRequest struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	Title   string `json:"title" gorm:"not null;size:50"`
	Content string `json:"content" gorm:"not null;size:100"`
}

type ListPostRequest struct {
	PageNo           int        `json:"page_no"`
	PageSize         int        `json:"page_size"`
	Title            string     `json:"title"`
	Content          string     `json:"content"`
	MinCommentNumber *uint      `json:"min_comment_number"`
	CreatedAtStart   *time.Time `json:"created_at_start"`
	CreatedAtEnd     *time.Time `json:"created_at_end"`
}

type PagePostResponse struct {
	Total int64
	Posts []Post
}

type ctxKey string

// 是否需要未每个业务单独配置一个变量 TODO
const ctxKeyId ctxKey = "ID"
const CtxKeyRowsAffected ctxKey = "RowsAffected"
const ctxKeyAudit ctxKey = "AuditBy"

type RowsAffectedContainer struct {
	Value  int64
	PostId uint
}

type CtxKeyAuditContext struct {
	AuditBy string
	UserId  uint
}

func ContextWithValue(id uint) context.Context {
	return context.WithValue(context.Background(), ctxKeyId, id)
}

func ContextValue(tx *gorm.DB) uint {
	if tx != nil && tx.Statement != nil && tx.Statement.Context != nil {
		if v, ok := tx.Statement.Context.Value(ctxKeyId).(uint); ok && v != 0 {
			return v
		}
	}
	return 0
}

func ContextWithValueAudit(auditBy *CtxKeyAuditContext) context.Context {
	return context.WithValue(context.Background(), ctxKeyAudit, auditBy)
}

func ContextValueAudit(tx *gorm.DB) *CtxKeyAuditContext {
	if tx != nil && tx.Statement != nil && tx.Statement.Context != nil {
		if v, ok := tx.Statement.Context.Value(ctxKeyAudit).(*CtxKeyAuditContext); ok && v != nil {
			return v
		}
	}
	return nil
}

func (a *Post) BeforeCreate(tx *gorm.DB) error {
	title := a.Title
	audit := ContextValueAudit(tx)
	if nil == audit {
		return nil
	}
	a.Audit.AuditBy = audit.AuditBy
	if utils.Audit.Check(title) {
		a.Audit.AuditStatus = "active"
	} else {
		a.Audit.AuditStatus = "inactive"
	}
	a.Audit.AuditVersion = "V1"
	// 对于嵌入字段，使用扁平的列名（snake_case）
	tx.Statement.SetColumn("audit_by", audit.AuditBy)
	tx.Statement.SetColumn("audit_status", a.Audit.AuditStatus)
	tx.Statement.SetColumn("audit_version", a.Audit.AuditVersion)
	return nil
}

// 文章创建钩子：计算用户文章数
// not execute when UNIQUE constraint failed: posts.title
func (a *Post) AfterCreate(tx *gorm.DB) error {
	fmt.Printf("文章创建钩子执行 AfterCreate\n")
	// userId := ContextValue(tx)
	audit := ContextValueAudit(tx)
	if nil == audit {
		return nil
	}
	fmt.Printf("文章创建钩子执行 AfterCreate 更新作者文章数+1 userId:%d", audit.UserId)
	var postUser User
	//在没有指定排序时，默认按主键升序排序
	//SELECT * FROM `users` WHERE `users`.`id` = 1 AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1
	if err := tx.First(&postUser, audit.UserId).Error; err != nil {
		return err
	}
	postUser.PostNumber++
	// UPDATE `users` SET `id`=1,`username`="admin",`email`="newadmin@example.com",`password`="$2a$10$eVi7VvODsGGh1Yr2g54OqOwNF7Z9ajHsKr8orXhxHtDriQT.iNmke",`post_number`=4,`created_at`="2026-01-14 11:46:31.684",`updated_at`="2026-01-15 12:36:56.566",`deleted_at`=NULL WHERE `users`.`deleted_at` IS NULL AND `id` = 1
	// if err := tx.Save(postUser).Error; err != nil {
	// UPDATE `users` SET `post_number`=5,`updated_at`="2026-01-15 12:39:35.354" WHERE `users`.`deleted_at` IS NULL AND `id` = 1
	if err := tx.Select("PostNumber").Save(postUser).Error; err != nil {
		return err
	}
	return nil
}
