package backend

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	pgConn *gorm.DB
)

func InitDB() {
	var dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("host"), os.Getenv("user"), os.Getenv("password"), os.Getenv("dbname"), os.Getenv("port"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	db.AutoMigrate(&Region{})
	pgConn = db
}

func allowedRegion(s string) bool {
	var (
		region Region
		tx     = pgConn.Where("name = ?", s).First(&region)
	)
	if tx.Error != nil {
		println(tx.Error.Error())
		return false
	}
	return region.Hash != ""
}

func insertNewRegion(name string) {
	var data, ok = latLongByName(name)
	if ok {
		pgConn.Create(&Region{
			Name: name,
			Lat:  data.Lat,
			Long: data.Long,
		})
	}
}

func updateHash(name, hash string) bool {
	if tx := pgConn.Model(&Region{}).Where("name = ?", name).Update("hash", hash); tx.Error != nil {
		println(tx.Error.Error())
		return false
	}
	return true
}

func getRegion(name string) (Region, bool) {
	var region Region
	if tx := pgConn.Where("name = ?", name).First(&region); tx.Error != nil {
		println(tx.Error.Error())
		return Region{}, false
	}
	return region, len(region.Hash) > 0
}
