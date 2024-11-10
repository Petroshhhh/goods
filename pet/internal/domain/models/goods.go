package models

type Goods struct {
	ID         int64  `db:"id"`
	Brand      string `db:"brand"`
	PlaceSave  int64  `db:"place_save"`
	StoreHouse int64  `db:"store_house"`
	Worker     string `db:"worker"`
}
