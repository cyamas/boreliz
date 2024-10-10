package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/cyamas/boreliz/internal/boulder"
	"github.com/cyamas/boreliz/internal/inventory"
)

func main() {
	args := os.Args

	inv := inventory.Load()
	if len(args) > 1 {
		switch args[1] {
		case "set":
			b := boulder.New()
			SetBoulder(b, inv)
		}
	}
	boulder := boulder.New()
	fmt.Println("boulder: ", boulder)
}

func SetBoulder(b *boulder.Boulder, inv *inventory.Inventory) {
	GetHolds(b, inv)
	SetMoves(b)
}

func GetHolds(b *boulder.Boulder, inv *inventory.Inventory) {
	fmt.Println("STEP 1: GET AND SET HOLDS")
	for {
		fmt.Println("ENTER HOLD ID to add to new boulder. Or enter 'I' to view the hold inventory:")
		fmt.Println("Enter 'F' to indicate that you are finished adding holds to your boulder")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Invalid entry.")
			continue
		}

		switch input[:len(input)-1] {
		case "F":
			return
		case "I":
			for _, hold := range inv.AllHolds() {
				fmt.Printf("ID: %d Manufacturer: %s Model: %s, Color: %s", hold.ID(), *hold.GetManufacturer(), hold.GetModel(), hold.GetColor())
			}
		default:
			id, err := strconv.Atoi(input)
			if err != nil {

			}
			hold, err := inventory.New().GetHoldByID(id)
			if err != nil {
				fmt.Println("hold id does not exist in inventory")
				continue
			}
			b.AddHold(hold)

		}
	}
}

func SetMoves(b *Boulder) {

	for {

	}
}
