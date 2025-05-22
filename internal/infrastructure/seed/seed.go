package seed

import (
	"amazing_gateway/internal/adapter/repository"
	"amazing_gateway/internal/infrastructure/database"
	"fmt"
	"github.com/bxcodec/faker/v3"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
)

func SeedUsers() {
	db := database.DB

	// Check si déjà seedé
	var count int64
	db.Model(&repository.User{}).Count(&count)
	if count > 0 {
		fmt.Println("Users already seeded")
		return
	}

	password, _ := bcrypt.GenerateFromPassword([]byte("hashed"), bcrypt.DefaultCost)

	// --- ADMINISTRATEURS ---
	for i := 0; i < 15; i++ {
		admin := repository.User{
			Name:     faker.Name(),
			Email:    fmt.Sprintf("admin%d@amazing.io", i),
			Password: string(password), // à adapter
			Role:     repository.Administrator,
		}
		db.Create(&admin)
	}

	// --- CLASSGROUPS (40 classes pour 5 années, 4 promos, 2 classes/promo) ---
	var classGroups []repository.ClassGroup
	for year := 1; year <= 5; year++ {
		for promo := 1; promo <= 4; promo++ {
			for i := 1; i <= 2; i++ {
				class := repository.ClassGroup{
					Name: fmt.Sprintf("Year %d Promo%d Class%d", year, promo, i),
				}
				db.Create(&class)
				classGroups = append(classGroups, class)
			}
		}
	}

	// --- TRAINEES ---
	for _, class := range classGroups {
		n := rand.Intn(6) + 15 // entre 15 et 20
		for i := 0; i < n; i++ {
			trainee := repository.User{
				Name:         faker.Name(),
				Email:        fmt.Sprintf("trainee_%d_%d@amazing.io", class.ID, i),
				Password:     string(password),
				Role:         repository.Trainee,
				ClassGroupID: &class.ID,
			}
			db.Create(&trainee)
		}
	}

	// --- TRAINERS ---
	// approx 40 trainers pour couvrir tous les cours à venir
	for i := 0; i < 40; i++ {
		trainer := repository.User{
			Name:     faker.Name(),
			Email:    fmt.Sprintf("trainer%d@amazing.io", i),
			Password: string(password),
			Role:     repository.Trainer,
		}
		db.Create(&trainer)
	}

	fmt.Println("✅ Users and class groups seeded")
}
