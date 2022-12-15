package categories

import (
	"github.com/itzngga/goRoxy/embed"
)

var CommonCategory = "Common"

func init() {
	embed.Categories.Add(CommonCategory)
}
