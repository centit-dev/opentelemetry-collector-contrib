package httpbodyprocessor

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Person struct {
	Id      *int    `json:"id"`
	Name    string  `json:"name"`
	Age     int     `json:"age"`
	Address Address `json:"address"`
}

type Address struct {
	City    string `json:"city"`
	Country string `json:"country"`
	Room    Room   `json:"room"`
}

type Room struct {
	Unit   string `json:"unit"`
	Number int    `json:"number"`
}

var person = Person{
	Id:   nil,
	Name: "John Doe",
	Age:  30,
	Address: Address{
		City:    "New York",
		Country: "USA",
		Room: Room{
			Unit:   "A",
			Number: 101,
		},
	},
}

func TestFlattenJson(t *testing.T) {
	jsonString, err := json.Marshal(person)
	assert.Nil(t, err)

	var jsonObject map[string]interface{}
	err = json.Unmarshal(jsonString, &jsonObject)
	assert.Nil(t, err)

	output := make(map[string]string)
	flattenJson(jsonObject, "", &output)

	_, ok := output["id"]
	assert.False(t, ok)
	assert.Equal(t, "John Doe", output["name"])
	assert.Equal(t, "30", output["age"])
	assert.Equal(t, "New York", output["address.city"])
	assert.Equal(t, "USA", output["address.country"])
	assert.Equal(t, "A", output["address.room.unit"])
	assert.Equal(t, "101", output["address.room.number"])
}

func TestProcessJson(t *testing.T) {
	jsonString := "{\"id\":\"1\",\"shopId\":\"1\",\"userId\":\"1\",\"avatarUrl\":\"http://www.baidu.com\",\"name\":\"shush\",\"gender\":\"1\",\"mobilePhone\":\"13770321976\",\"level\":null,\"birthday\":null,\"country\":null,\"province\":null,\"city\":null,\"district\":null,\"labels\":null,\"memo\":null}{\"id\":\"1\",\"shopId\":\"1\",\"userId\":\"1\",\"avatarUrl\":\"http://www.baidu.com\",\"name\":\"shush\",\"gender\":\"1\",\"mobilePhone\":\"13770321976\",\"level\":null,\"birthday\":null,\"country\":null,\"province\":null,\"city\":null,\"district\":null,\"labels\":null,\"memo\":null}"
	object := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonString), &object)
	assert.Nil(t, err)
}
