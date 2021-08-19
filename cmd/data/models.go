package data

import (
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type Models struct {
	Keepers KeeperModel
	Bees    BeeModel

}

func NewModels(db *gorm.DB) Models {
	return Models{
		Keepers: KeeperModel{DB: db},
		Bees:    BeeModel{DB: db},
	}
}

type KeeperModel struct {
	DB *gorm.DB
}

type BeeModel struct {
	DB *gorm.DB
}

type Token struct {
	UserID uint
	jwt.StandardClaims
}
