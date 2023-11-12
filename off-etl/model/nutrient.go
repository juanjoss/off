package model

type NutrientLevels struct {
	Fat          string `json:"fat" db:"fat"`
	SaturatedFat string `json:"saturated-fat" db:"saturated_fat"`
	Sugar        string `json:"sugars" db:"sugar"`
	Salt         string `json:"salt" db:"salt"`
}

type Nutriments struct {
	Energy100g    float32 `json:"energy-kcal_100g"`
	EnergyServing float32 `json:"energy-kcal_serving"`
	NOVA          uint8   `json:"nova-group" db:"nova_group"`
}

type NutriscoreData struct {
	Score int8   `json:"score" db:"nutriscore"`
	Grade string `json:"grade" db:"nutriscore_grade"`
}
