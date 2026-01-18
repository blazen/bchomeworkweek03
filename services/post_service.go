package services

import (
	"gorm.io/gorm"

	"gin-examples/project/models"
	"gin-examples/project/utils"
)

type PostService struct {
	db *gorm.DB
}

func NewPostService(db *gorm.DB) *PostService {
	return &PostService{db: db}
}

// 创建文章
func (s *PostService) CreatePost(userId uint, req models.CreatePostRequest) (*models.Post, error) {
	// 检查标题是否已存在
	var existingPost models.Post
	if err := s.db.Where("title = ?", req.Title).First(&existingPost).Error; err == nil {
		return nil, utils.NewAppError(409, "Post title exist")
	}

	// ctx := models.ContextWithValue(userId)
	// ctx = models.ContextWithValueAudit("blob")
	container := &models.CtxKeyAuditContext{
		AuditBy: "bob",
		UserId:  userId,
	}
	ctx := models.ContextWithValueAudit(container)

	// 创建文章
	post := models.Post{
		UserID:  userId,
		Title:   req.Title,
		Content: req.Content,
	}

	if err := s.db.WithContext(ctx).Create(&post).Error; err != nil {
		return nil, err
	}

	return &post, nil
}

// 查询用户的全部文章
func (s *PostService) ListPost(userId uint, pageNo, pageSize int) ([]models.Post, error) {
	var posts []models.Post
	if err := s.db.Scopes(utils.Sql.Paginate(pageNo, pageSize), utils.Sql.OrderCreateAt()).
		Where("user_id = ?", userId).
		// Order("created_at desc").
		Find(&posts).Error; err != nil {
		return nil, utils.NewAppError(409, "Query Post failed by userId")
	}
	return posts, nil
}

// 查询所有用户的文章
func (s *PostService) ListPostAll(pageNo, pageSize int) (*models.PagePostResponse, error) {
	var posts []models.Post
	if err := s.db.Scopes(utils.Sql.Paginate(pageNo, pageSize), utils.Sql.OrderCreateAt()).
		Where("audit_status", "active"). // 进查询 title 审计通过的
		// Order("created_at desc").
		Find(&posts).Error; err != nil {
		return nil, utils.NewAppError(409, "Query Post failed")
	}
	var total int64
	if err := s.db.Model(&models.Post{}).
		Where("audit_status", "active"). // 进查询 title 审计通过的
		Count(&total).Error; err != nil {
		return nil, utils.NewAppError(409, "Query Post failed")
	}
	page := &models.PagePostResponse{
		Total: total,
		Posts: posts,
	}
	return page, nil
}

// 条件查询文章
func (s *PostService) ListPostByCondition(req models.ListPostRequest) ([]models.Post, error) {
	var posts []models.Post
	tx := s.db
	// 动态拼接条件：创建时间范围（仅当Start和End都传了才筛选）
	if req.CreatedAtStart != nil && req.CreatedAtEnd != nil {
		tx = tx.Where("created_at BETWEEN ? AND ?", req.CreatedAtStart, req.CreatedAtEnd)
	}
	// 动态拼接条件：最小评论数
	if req.MinCommentNumber == nil {
		tx = tx.Where("comment_number >= ?", 0)
	} else {
		tx = tx.Where("comment_number >= ?", req.MinCommentNumber)
	}
	// 动态拼接条件：标题包含
	if req.Title != "" {
		tx = tx.Where("title like ?", "%"+req.Title+"%")
	}
	tx = tx.Where("audit_status", "active") // 进查询 title 审计通过的
	// SELECT * FROM `posts` WHERE comment_number >= 0 AND title >= "%hello%" AND `posts`.`deleted_at` IS NULL ORDER BY created_at desc LIMIT 10
	// SELECT * FROM `posts` WHERE comment_number >= 2 AND title >= "%hello%" AND `posts`.`deleted_at` IS NULL ORDER BY created_at desc LIMIT 10
	// SELECT * FROM `posts` WHERE (created_at BETWEEN "2026-01-10 12:39:35.35" AND "2026-01-16 15:05:28.322") AND comment_number >= 2 AND title >= "%hello%" AND `posts`.`deleted_at` IS NULL ORDER BY created_at desc LIMIT 10
	if err := tx.Scopes(utils.Sql.Paginate(req.PageNo, req.PageSize), utils.Sql.OrderCreateAt()).
		// Where("title = ?", req.Title).
		// Order("created_at desc").
		Find(&posts).Error; err != nil {
		return nil, utils.NewAppError(409, "Query Post failed by condition")
	}
	return posts, nil
}

// 主键查询文章，关联查询最新两条评论
func (s *PostService) GetPostById(id uint) (*models.Post, error) {
	// First with primary key
	var post models.Post
	// SELECT * FROM `posts` WHERE `posts`.`id` = 11 AND `posts`.`deleted_at` IS NULL ORDER BY `posts`.`id` LIMIT 1
	// SELECT * FROM `comments` WHERE `comments`.`post_id` = 11 AND `comments`.`deleted_at` IS NULL ORDER BY created_at desc LIMIT 2
	if err := s.db.Preload("Comments",
		func(db *gorm.DB) *gorm.DB {
			// return db.Scopes(utils.Sql.Paginate(1, 2)).Order("created_at desc")
			return db.Scopes(utils.Sql.Paginate(1, 2), utils.Sql.OrderCreateAt())
		}).First(&post, id).Error; err != nil {
		return nil, utils.NewAppError(409, "Query Post failed by id")
	}
	return &post, nil
}

// 查询评论数量最多的文章，关联查询最新两条评论
func (s *PostService) GetPostByMaxCommentNumber() (*models.Post, error) {
	var post models.Post
	// SELECT * FROM `posts` WHERE `posts`.`deleted_at` IS NULL ORDER BY comment_number desc,`posts`.`id` LIMIT 1
	if err := s.db.Preload("Comments",
		func(db *gorm.DB) *gorm.DB {
			// return db.Scopes(utils.Sql.Paginate(1, 2)).Order("created_at desc")
			return db.Scopes(utils.Sql.Paginate(1, 2), utils.Sql.OrderCreateAt())
		}).Order("comment_number desc").First(&post).Error; err != nil {
		return nil, utils.NewAppError(409, "Query max CommentNumber of Post failed")
	}
	return &post, nil
}

// 更新用户的文章
func (s *PostService) UpdatePost(userId uint, req models.UpdatePostRequest) (*models.Post, error) {
	// 检查标题是否已存在
	var existingPost models.Post
	if err := s.db.Where("id = ? and user_id = ?", req.ID, userId).First(&existingPost).Error; err != nil {
		return nil, utils.NewAppError(409, "Post not exist")
	}

	existingPost.Title = req.Title
	existingPost.Content = req.Content

	if err := s.db.Save(existingPost).Error; err != nil {
		return nil, err
	}

	return &existingPost, nil
}

// 删除用户的文章
func (s *PostService) DeletePost(userId uint, id uint) (bool, error) {
	// 检查标题是否已存在
	var existingPost models.Post
	// physical delete by unscoped
	// if err := s.db.Where("id = ? and user_id = ?", id, userId).Delete(&existingPost).Error; err != nil {
	if err := s.db.Where("id = ? and user_id = ?", id, userId).Unscoped().Delete(&existingPost).Error; err != nil {
		return false, utils.NewAppError(409, "Post delete failed")
	}
	return true, nil
}
