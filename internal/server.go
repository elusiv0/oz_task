package main

import (
	"fmt"

	"github.com/Masterminds/squirrel"
	model "github.com/elusiv0/oz_task/internal/dto"
)

const defaultPort = "8080"

func main() {
	some := 1
	some2 := 3
	some4 := 5
	commentsReq := model.GetCommentsRequest{
		ParentId: &some,
		PostId:   &some4,
		First:    some2,
	}
	doThis(commentsReq)
}

func doThis(req model.GetCommentsRequest) {
	var conditions squirrel.And
	if req.ParentId != nil {
		conditions = append(conditions, squirrel.Eq{"parent_id": *req.ParentId})
	}
	if req.PostId != nil {
		conditions = append(conditions, squirrel.Eq{"article_id": *req.PostId})
	}
	if req.After != nil {
		conditions = append(conditions, squirrel.Lt{"id": *req.After})
	}
	fmt.Println(conditions.ToSql())
	fmt.Println(req.ParentId)
	fmt.Println(*req.ParentId)
	fmt.Println(*req.After)
}
