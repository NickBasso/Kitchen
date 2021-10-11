package kitchen

import (
	"kitchen/src/cook"
	"kitchen/src/cookingApparatus"
)


type Kitchen struct {
	cooks []cook.Cook
	cookingApparatus []cookingApparatus.CookingApparatus
}
