package data

import (
	"gorm.io/gorm"
	"strconv"
)

type caste uint8

type Bee struct {
	gorm.Model
	Name       string `json:"name"`
	Species    string `json:"species"`
	Caste      caste  `json:"caste"`
	KeeperID   uint   `json:"keeper_id"`
}

func (c caste) MarshalJSON() ([]byte, error) {
	var t string
	switch c {
	case 0:
		t = "Worker Bee"
	case 1:
		t = "Drone Bee"
	case 2:
		t = "The Queen Bee"
	default:
		t = "Worker Bee"
	}
	quoted := strconv.Quote(t)
	return []byte(quoted), nil
}



func (c caste) String() string {
	switch c {
	case 0:
		return "Worker"
	case 1:
		return "Drone"
	case 2:
		return "Queen"
	default:
		return "Little One"
	}
}


func (m *BeeModel) Insert(bee *Bee) error {
	err := m.DB.Create(bee)

	if err != nil {
		return err.Error
	}

	if err.RowsAffected != 1{
		return err.Error
	}
	return nil
}

func (m *BeeModel) GetAll(keeperID uint) ([]*Bee, error) {
	var bees []*Bee
	err := m.DB.Table("bees").Where("keeper_id = ?", keeperID).Find(&bees).Error
	if err != nil {
		return nil, err
	}

	return bees, nil
}
