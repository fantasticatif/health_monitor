package endpointupdater

import (
	"fmt"
	"log"

	"github.com/fantasticatif/health_monitor/api/db"
	"github.com/fantasticatif/health_monitor/api/shareddata"
	"github.com/fantasticatif/health_monitor/data"
	"github.com/fantasticatif/health_monitor/util"
	"github.com/joho/godotenv"
)

func getAccounts() *[]data.Account {
	var accounts []data.Account
	tx := db.SharedDB.Where("deleted_at", nil).Find(&accounts)
	if tx.Error != nil {
		log.Fatalf("error fetching accounts %s\n", tx.Error.Error())
	}
	return &accounts
}

func processAccount(acc data.Account) {
	log.Printf("--- processing acc %s begin\n", acc.Name)
	var hps []data.HitPoint
	tx := db.SharedDB.Where("account_id", acc.ID).Find(&hps)
	if tx.Error != nil {
		log.Printf("error getting hitpoints for account %d : %s", acc.ID, tx.Error.Error())
	} else {
		log.Printf("found account %d hp's %d", acc.ID, len(hps))
		for _, hp := range hps {
			current_status, err := hp.GetStatus()
			if err != nil {
				continue
			}
			if current_status != hp.Status {
				fmt.Printf("hp status changed %s from %s\n", current_status, hp.Status)
				tx := db.SharedDB.Model(&hp).Where("id = ?", hp.ID).Update("status", current_status)
				if tx.Error != nil {
					log.Printf("update hitpoint failed %s, hp id: %d, err: %s\n", tx.Error.Error(), hp.ID, tx.Error.Error())
				} else {
					log.Printf("updated hitpoint hp id: %d\n", hp.ID)
				}
			}
		}
	}
}

func Run() {
	log.Println("update endpoints")
	godotenv.Load()
	shareddata.ResetEnvVariables()
	db.SetupDb()

	accounts := getAccounts()

	// Need to be initialized by initialize method on creation before use
	coordinator := util.RoutineCoordinator{}
	coordinator.Initialize(8)

	for _, acc := range *accounts {
		processAcc := func(vals ...any) {
			processAccount(acc)
		}
		coordinator.Add(processAcc)
		log.Printf("account %s", acc.Name)
	}

	coordinator.WaitForCompletion()

}
