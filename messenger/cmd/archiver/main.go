package main

import (
	"fmt"
	"log"

	ssov5 "github.com/mgrankin-cloud/messenger/contract/gen/go/user"
	"github.com/mgrankin-cloud/messenger/pkg/utils/lib/archiver"
)

func main() {
	userID, err := ssov5.GetUserResponse{
		UserId: userID,
	}
	if err != nil {
		log.Fatal(err)
	}
	files := []string{req.GetFiles()}

	err = archiver.SendArchive(nil, userID, files)
	if err != nil {
		fmt.Printf("Ошибка при отправке архива: %v\n", err)
		return
	}

	fmt.Printf("Архив отправлен пользователю %d\n", userID)
}
