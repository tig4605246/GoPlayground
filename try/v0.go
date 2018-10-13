package her

import (
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

type GameInfo struct {
	ID        bson.ObjectId `json:"ID,omitempty" bson:"_id"`
	Name      string        `json:"name,omitempty" bson:"name,omitempty"`
	Developer string        `json:"develop,omitempty" bson:"develop,omitempty"`
	Type      string        `json:"type,omitempty" bson:"type,omitempty"`
	Platform  string        `json:"platform,omitempty" bson:"platform,omitempty"`
	Rated     string        `json:"rated,omitempty" bson:"rated,omitempty"`
	Sale      int           `json:"sale,omitempty" bson:"sale,omitempty"`
}

type FormInfo struct {
	Name       string
	Type       string
	UpperBound int
	LowerBound int
}

func GameInfoform() *GameInfo {
	return &GameInfo{}
}

const (
	test = "testdb"
	_3A  = "AAA_level_games"
)

func main() {

	fmt.Printf("hello, world\n")
	fmt.Println("Going to test the validator\n")
	x := "good work"

	if govalidator.Contains(x, "good") {
		fmt.Println("Got it")
	}
	fmt.Println()

	inputInfo := CreateInputMap()
	name := [6]string{"Name","Developer","Type","Platform","Rated","Sale"}
	vType := [6]string{"Alphanumeric", "Alphanumeric", "Alphanumeric", "Alphanumeric", "Alphanumeric", "Numeric"}

	map[string]string {
		"user": "Alphanumeric,email,range(1:4),"}
	}
	targetInfo := CreateValidationMap(name,vType)

	res, ok := CheckParam(inputInfo, targetInfo)
	fmt.Println()
	res, ok = CheckRange(inputInfo, targetInfo)
	if res != true {
		fmt.Println("Validation failed, ", ok)
	} else {

		fmt.Println(inputInfo, "\n", targetInfo)
	}

	return

	//Connect to  local db

	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	game := GenerateGameInfo(true)

	c := session.DB(test).C(_3A)

	err = c.Insert(game)
	if err != nil {
		fmt.Println(err)
	}

}

func GenerateGameInfo(isLegal bool) *GameInfo {
	newGame := GameInfoform()
	newGame.ID = bson.NewObjectId()
	newGame.Name = "TheElderScrollVSkyrim"
	newGame.Developer = "Dunkey"
	newGame.Type = "RPG,FPS"
	newGame.Platform = "PS4,PS3,PC,Xbox,XboxOne,NS"
	newGame.Rated = "M"
	newGame.Sale = 5000000

	return newGame
}

func CreateInputMap() map[string]string {
	testInfo := map[string]string{
		"Name":      "FalloutNewVegas",
		"Developer": "Bethesda",
		"Type":      "RPG",
		"Platform":  "PC",
		"Rated":     "M",
		"Sale":      "1000",
	}
	return testInfo
}

func CreateValidationMap(name [6]string, vType [6]string) map[int]FormInfo {

	validInfo := make(map[int]FormInfo)
	newGame := FormInfo{}

	if len(name) != len(vType) {
		fmt.Println("param length incorrect")
		return nil
	}
 
	for i := 0 ; i < len(name) ; i++{
		newGame.Name = name[i]
		newGame.Type = vType[i]
		newGame.LowerBound = 500
		newGame.UpperBound = 700
		validInfo[i] = newGame
	}
	return validInfo
}

