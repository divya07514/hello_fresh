package model

type RecipeData struct {
	Postcode string `json:"postcode"`
	Recipe   string `json:"recipe"`
	Delivery string `json:"delivery"`
}

type RecipeStats struct {
	UniqueRecipes           int                   `json:"unique_recipe_count"`
	PerRecipe               []*PerRecipeStats     `json:"count_per_recipe"`
	BusiestPostCode         *PostCodeStats        `json:"busiest_postcode"`
	CountPerPostCodeAndTime *PostCodeAndTimeStats `json:"count_per_postcode_and_time"`
	MatchByName             []string              `json:"match_by_name"`
}

type PerRecipeStats struct {
	Recipe string `json:"recipe"`
	Count  int64  `json:"count"`
}
type PostCodeStats struct {
	Postcode      string `json:"postcode"`
	DeliveryCount int64  `json:"delivery_count"`
}

type PostCodeAndTimeStats struct {
	PostCode      string `json:"postcode"`
	FromTime      string `json:"from"`
	ToTime        string `json:"to"`
	DeliveryCount int64  `json:"delivery_count"`
}
