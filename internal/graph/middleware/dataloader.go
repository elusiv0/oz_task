package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/elusiv0/oz_task/internal/dto"
	"github.com/elusiv0/oz_task/internal/graph/dataloader"
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
				for idx, val := range commentReqs {
					insertIdx := 0
					if val.ParentId != nil {
						insertIdx = *val.ParentId
					} else {
						insertIdx = *val.PostId
					}
					mapReqIndex[insertIdx] = idx
				}
				comments := make(map[int]bool)
				commentsDto := make([][]*dto.Comment, len(commentReqs))

				commentResp, err := s.GetMany(r.Context(), commentReqs...)
				if err != nil {
					return [][]*dto.Comment{}, []error{err}
				}
				last := 0
				prev := -1
				for idx := 0; idx < len(commentResp); idx++ {
					val := commentResp[idx]
					if commentReqs[0].ParentId != nil {
						_, ok := comments[*val.ParentID]
						if !ok || idx+1 == len(commentResp) {
							insertIdx := idx
							if idx+1 == len(commentResp) {
								insertIdx++
							}

							if insertIdx != 0 {
								prev = *commentResp[insertIdx-1].ParentID
								commentsDto[mapReqIndex[prev]] = commentResp[last:insertIdx]
								last = idx
							}
						}
						if ok {
							continue
						}
						comments[*val.ParentID] = true
					} else {
						_, ok := comments[val.ArticleID]
						if !ok || idx+1 == len(commentResp) {
							insertIdx := idx
							if idx+1 == len(commentResp) {
								insertIdx++
							}

							if insertIdx != 0 {
								prev = commentResp[insertIdx-1].ArticleID
								commentsDto[mapReqIndex[prev]] = commentResp[last:insertIdx]
								last = idx
							}
						}
						if ok {
							continue
						}
						comments[val.ArticleID] = true
					}
				}
				errorsResp := make([]error, len(commentReqs))

				return commentsDto, errorsResp
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
