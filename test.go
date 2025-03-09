package ToDoGo

import (
	"fmt"
)

type TodoItem1 struct { //Структура члена списка
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

func main() {
	oldItem := TodoItem1{
		Id:          1,
		Title:       "Test",
		Description: "Test description",
		Done:        false,
	}

	newItem := TodoItem1{
		Id:          1,
		Title:       "Test",
		Description: "Test description",
		Done:        false,
	}
	var ResItem TodoItem1
	if newItem.Done != oldItem.Done && &newItem.Done != nil {
		ResItem.Done = newItem.Done
		fmt.Println("done:", newItem.Done)
	} else {
		fmt.Println("missing done")
	}
}
