package reqhandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"secondMS/repository"
	"secondMS/repository/models"
)

func AddPerson(rawData []byte) {
	ctx := context.Background()
	var person models.Person
	err := json.Unmarshal(rawData, &person)
	if err != nil {
		fmt.Println("Error unmarshaling data")
		return
	}
	err = repository.AddressBookDb.StoreOnePerson(ctx, person)
	if err != nil {
		fmt.Println("Error saving person, error: " + err.Error())
	}
	fmt.Println("Person added correctly")
}
