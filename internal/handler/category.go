package handler

type Category struct {
	Name        string
	Description string
}

var (
	UtilitiesCategory = &Category{
		Name:        "Utilities",
		Description: "Bot Utilities",
	}
	MiscCategory = &Category{
		Name:        "Misc",
		Description: "Bot Misc",
	}
	Uncategorized = &Category{
		Name:        "Uncategorized",
		Description: "Bot Uncategorized",
	}
)
