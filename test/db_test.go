package test

import (
	"inssa_club_waitlist_backend/cmd/server/models"
	"inssa_club_waitlist_backend/cmd/server/utils"
	"testing"

	_ "github.com/joho/godotenv/autoload"
	"github.com/tj/assert"
	"gopkg.in/guregu/null.v4"
)

func TestDBSingleton(t *testing.T) {
	db1 := utils.GetDB()
	db2 := utils.GetDB()
	utils.GetDB().SetupDB()

	assert.Equal(t, true, db1.Instance == db2.Instance)
}

func TestDBSetup(t *testing.T) {
	db := utils.GetDB()
	db.SetupDB()
	db.Instance.Exec("DELETE FROM interests") // reset db data
}

func TestRecordCreation(t *testing.T) {
	CLUBHOUSE_USER_ID := null.NewInt(1234, true)
	const EMAIL = "code.yeon.gyu@gmail.com"
	db := utils.GetDB()
	db.SetupDB()

	interest := models.Interest{ClubhouseUserID: CLUBHOUSE_USER_ID, Email: EMAIL}
	db.Instance.Create(&interest)

	result := models.Interest{}
	db.Instance.Where(models.Interest{ClubhouseUserID: CLUBHOUSE_USER_ID}).First(&result)
	// retrieve the same data

	assert.Equal(t, true, interest.Email == result.Email)
	// check if the requested data and created data is the same
}
