package db

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/muhwyndhamhp/todo-mx/config"
	"github.com/muhwyndhamhp/todo-mx/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var database *gorm.DB

func init() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
		config.Get(config.DB_HOST),
		config.Get(config.DB_PORT),
		config.Get(config.DB_USER),
		config.Get(config.DB_NAME),
		config.Get(config.DB_PASSWORD),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	database = db
}

func GetDB() *gorm.DB {
	return database
}

func SeedDB() error {
	todos := []models.Todo{
		{
			Title: "Clean the house!",
			Body: pgtype.Text{
				String: "You need to wipe the floor, mop it, and clean the tables",
				Valid:  true,
			},
			EncodedBody: "You need to wipe the floor, mop it, and clean the tables",
		},
		{
			Title: "Kiss Ma Wife!",
			Body: pgtype.Text{
				String: "kiss your wife for happinex x100",
				Valid:  true,
			},
			EncodedBody: "kiss your wife for happinex x100",
		},
		{
			Title: "The Thing...",
			Body: pgtype.Text{
				String: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum tincidunt erat vulputate vehicula gravida. Nullam tincidunt vehicula lorem ac ultricies. Proin elit libero, dignissim sed dolor sed, aliquet euismod velit. Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Praesent fringilla luctus consequat. In hac habitasse platea dictumst. Duis efficitur purus ante, sed pretium nibh egestas eget. Curabitur et viverra orci, venenatis rhoncus sapien. Sed varius mattis elit, sit amet sollicitudin turpis vestibulum a. Ut quam leo, lobortis quis maximus quis, blandit eget mi. Duis nisi massa, dictum ut faucibus eu, mollis ac ipsum. Pellentesque tristique id diam et mollis. Curabitur accumsan ipsum nec turpis laoreet, at tincidunt elit euismod. Vivamus ante erat, porttitor id lacus ac, eleifend efficitur mi. Etiam molestie nisl in mollis porttitor. Maecenas ligula eros, placerat vel facilisis eget, varius in dolor. ",
				Valid:  true,
			},
			EncodedBody: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum tincidunt erat vulputate vehicula gravida. Nullam tincidunt vehicula lorem ac ultricies. Proin elit libero, dignissim sed dolor sed, aliquet euismod velit. Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Praesent fringilla luctus consequat. In hac habitasse platea dictumst. Duis efficitur purus ante, sed pretium nibh egestas eget. Curabitur et viverra orci, venenatis rhoncus sapien. Sed varius mattis elit, sit amet sollicitudin turpis vestibulum a. Ut quam leo, lobortis quis maximus quis, blandit eget mi. Duis nisi massa, dictum ut faucibus eu, mollis ac ipsum. Pellentesque tristique id diam et mollis. Curabitur accumsan ipsum nec turpis laoreet, at tincidunt elit euismod. Vivamus ante erat, porttitor id lacus ac, eleifend efficitur mi. Etiam molestie nisl in mollis porttitor. Maecenas ligula eros, placerat vel facilisis eget, varius in dolor. ",
		},
	}

	return GetDB().Save(&todos).Error
}
