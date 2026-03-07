package pagination

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

const defaultLimit = 20
const maxLimit = 100

// Params holds pagination query parameters.
type Params struct {
	Page   int
	Limit  int
	Offset int
}

// Parse extracts ?page=&limit= from the Gin context.
func Parse(c *gin.Context) Params {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", fmt.Sprint(defaultLimit)))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	return Params{
		Page:   page,
		Limit:  limit,
		Offset: (page - 1) * limit,
	}
}

// SetTotalCount writes the X-Total-Count response header.
func SetTotalCount(c *gin.Context, total int64) {
	c.Header("X-Total-Count", fmt.Sprint(total))
}
