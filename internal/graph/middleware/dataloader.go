package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/elusiv0/oz_task/internal/dto"
	"github.com/elusiv0/oz_task/internal/graph/dataloader"
	repo "github.com/elusiv0/oz_task/internal/repo"
	"github.com/elusiv0/oz_task/internal/service"
)

const (
	commentLoaderKey = "commentLoader"
)

func DataloaderMiddleware(s service.CommentService, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		commentLoaderConfig := dataloader.CommentLoaderConfig{
			MaxBatch: 100,
			Wait:     5 * time.Millisecond,
			Fetch: func(commentReqs []dto.GetCommentsRequest) ([][]*dto.Comment, []error) {
				mapReqIndex := make(map[int]int)
				byParent := false
				if commentReqs[0].ParentId != nil {
					byParent = true
				}
				for idx, val := range commentReqs {
					insertIdx := 0
					if byParent {
						insertIdx = *val.ParentId
					} else {
						insertIdx = *val.PostId
					}
					mapReqIndex[insertIdx] = idx
				}
				comments := make(map[int]bool)
				errorsResp := make([]error, len(commentReqs))
				commentsResp := make([][]*dto.Comment, len(commentReqs))

				commentsPool, err := s.GetMany(r.Context(), commentReqs...)
				if err != nil {
					for idx, _ := range errorsResp {
						errorsResp[idx] = err
					}
					return commentsResp, errorsResp
				}
				last := 0
				prev := -1
				cplen := len(commentsPool)

				for idx := 0; idx < cplen; idx++ {
					val := commentsPool[idx]
					if byParent {
						_, ok := comments[*val.ParentID]
						if !ok && idx != 0 {
							insertIdx := idx

							prev = *(commentsPool[insertIdx-1].ParentID)
							commentsResp[mapReqIndex[prev]] = commentsPool[last:insertIdx]
							last = idx
						}
						if idx+1 == cplen {
							commentsResp[mapReqIndex[*val.ParentID]] = commentsPool[last:]
						}
						if ok {
							continue
						}
						comments[*val.ParentID] = true
					} else {
						_, ok := comments[val.ArticleID]
						if !ok && idx != 0 {
							insertIdx := idx
							if idx+1 == cplen {
								insertIdx++
							}

							prev = commentsPool[insertIdx-1].ArticleID
							commentsResp[mapReqIndex[prev]] = commentsPool[last:insertIdx]
							last = idx
						}
						if idx+1 == cplen {
							commentsResp[mapReqIndex[val.ArticleID]] = commentsPool[last:]
						}
						if ok {
							continue
						}
						comments[val.ArticleID] = true
					}
				}

				for idx, val := range commentsResp {
					if len(val) == 0 {
						errorsResp[idx] = dto.NewCustomError(repo.CommentsNotFoundErr, commentReqs[idx])
					}
				}
				return commentsResp, errorsResp
			},
		}
		commentLoader := dataloader.NewCommentLoader(commentLoaderConfig)

		ctx := context.WithValue(r.Context(), commentLoaderKey, commentLoader)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetCommentLoader(ctx context.Context) *dataloader.CommentLoader {
	return ctx.Value(commentLoaderKey).(*dataloader.CommentLoader)
}
