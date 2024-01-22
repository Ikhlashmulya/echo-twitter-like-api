package main

import (
	"flag"
	"os"

	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/config"
	"github.com/Ikhlashmulya/echo-twitter-like-api/internal/entity"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var configuration *viper.Viper
var log *logrus.Logger
var db *gorm.DB

func init() {
	configuration = config.NewViper()
	log = config.NewLogger(configuration)
	db = config.NewGorm(configuration, log)
}

func main() {
	command := flag.String("migrate", "up", "migration opration")
	flag.Parse()

	migrate := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "202401220815",
			Migrate: func(db *gorm.DB) error {
				return db.Migrator().AutoMigrate(&entity.User{})
			},
			Rollback: func(db *gorm.DB) error {
				return db.Migrator().DropTable("users")
			},
		},
		{
			ID: "202401220821",
			Migrate: func(db *gorm.DB) (err error) {
				if err = db.Migrator().AutoMigrate(&entity.Post{}); err != nil {
					return
				}

				if err = db.Migrator().CreateConstraint(&entity.Post{}, "Author"); err != nil {
					return
				}

				return
			},
			Rollback: func(db *gorm.DB) (err error) {
				return db.Migrator().DropTable("posts")
			},
		},
		{
			ID: "202401220830",
			Migrate: func(db *gorm.DB) (err error) {
				if err = db.Migrator().AutoMigrate(&entity.Comment{}); err != nil {
					return
				}

				if err = db.Migrator().CreateConstraint(&entity.Comment{}, "Post"); err != nil {
					return
				}

				if err = db.Migrator().CreateConstraint(&entity.Comment{}, "User"); err != nil {
					return
				}

				return

			},
			Rollback: func(db *gorm.DB) (err error) {
				return db.Migrator().DropTable("comments")
			},
		},
	})

	var err error

	if *command == "up" {
		if version := isSetVersion(os.Args); version != "" {
			err = migrate.MigrateTo(version)
		} else {
			err = migrate.Migrate()
		}
	} else if *command == "down" {
		if version := isSetVersion(os.Args); version != "" {
			err = migrate.RollbackTo(version)
		} else {
			err = migrate.RollbackLast()
		}
	} else {
		log.Fatalf("invalid command")
	}

	if err != nil {
		log.Fatalf("error database migration: %v", err)
	} else {
		log.Println("migration success")
	}
}

func isSetVersion(args []string) string {
	if len(args) > 3 {
		return os.Args[3]
	}

	return ""
}
