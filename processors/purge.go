package processors

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/swims/nadeta-sql/migrations"
)

func RunPurge(migrationStore *migrations.MigrationStore) error {
	fmt.Print("Are you sure you want to purge all data stored? (Y/y)\n")
	reader := bufio.NewReader(os.Stdin)
	userInput, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	userInput = strings.TrimSpace(userInput)
	if userInput == "y" || userInput == "Y" {
		err := migrationStore.PurgeData()
		if err != nil {
			return err
		}
		fmt.Println("Removed all data from the system...")
	} else {
		fmt.Println("Cancelling Purge...")
	}
	return nil
}
