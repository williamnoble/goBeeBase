package data

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"os"
)

var (
	ErrRecordNotFound = errors.New("the specified user could not be found")
)

type Keeper struct {
	gorm.Model
	Email      string `json:"email,omitempty" json:"email,omitempty"`
	Password   string `json:"password,omitempty" json:"password,omitempty"`
	Token      string `json:"token";sql:"-"`
}

type KeeperResponse struct {
	KeeperID    uint
	KeeperEmail string
	Token       string
}

func (m *KeeperModel) Create(keeper *Keeper) (*KeeperResponse, error) {
	var pass struct {
		plaintextPassword string
		hashedPassword    string
	}

	pass.plaintextPassword = keeper.Password

	hash, err := bcrypt.GenerateFromPassword([]byte(pass.plaintextPassword), 12)
	if err != nil {
		return nil, err
	}

	pass.hashedPassword = string(hash)
	keeper.Password = pass.hashedPassword

	claim := &Token{UserID: keeper.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claim)
	tokenToString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	keeper.Token = tokenToString

	m.DB.Create(keeper)
	k := KeeperResponse{
		KeeperID:    keeper.ID,
		//KeeperID:    uint(keeper.ID),
		KeeperEmail: keeper.Email,
		Token:       keeper.Token,
	}
	return &k, nil

}

func (m *KeeperModel) GetByEmail(email string) (Keeper, error) {
	k := &Keeper{}
	err := m.DB.Table("keepers").Where("email = ?", email).First(k).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return  Keeper{}, ErrRecordNotFound
		default:
			return Keeper{}, err
		}
	}
	return *k, nil
}
