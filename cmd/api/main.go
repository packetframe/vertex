package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"

	"github.com/packetframe/vertex/internal/config"
	"github.com/packetframe/vertex/internal/db"
)

var version = "dev"

var (
	dbDsn     = os.Getenv("DB_DSN")
	sentryDsn = os.Getenv("SENTRY_DSN")
)

// response returns a JSON response
func response(c *fiber.Ctx, status int, message string, data map[string]interface{}) error {
	// Capitalize first letter
	if len(message) > 1 {
		message = strings.ToUpper(message[0:1]) + message[1:]
	}

	return c.Status(status).JSON(fiber.Map{
		"success": (200 <= status) && (status < 300),
		"message": message,
		"data":    data,
	})
}

// internalServerError logs and returns an internal server error
func internalServerError(c *fiber.Ctx, err error) error {
	fmt.Printf("Internal Server Error ---------------------- %s ----------------------\n", err)
	sentry.CaptureException(err)
	return response(c, http.StatusInternalServerError, "Internal Server Error", nil)
}

func main() {
	if version == "dev" {
		log.SetLevel(log.DebugLevel)
	}

	if sentryDsn != "" {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:     sentryDsn,
			Release: version,
		}); err != nil {
			log.Fatalf("sentry.Init: %s", err)
		}
	}

	if dbDsn == "" {
		log.Warn("DB_DSN is not set, using development database")
		dbDsn = "host=localhost user=api password=api dbname=api port=5432 sslmode=disable"
	}

	log.Infof("Connecting to database")
	database, err := db.Connect(dbDsn)
	if err != nil {
		log.Fatal(err)
	}

	metricsTicker := time.NewTicker(2 * time.Second)
	go func() {
		for range metricsTicker.C {
			log.Debugln("Looking for expired rules")

			var rules []db.Rule
			if err := database.Find(&rules).Error; err != nil {
				log.Warnf("unable to retreive rules: %s", err)
			}

			for _, rule := range rules {
				if time.Since(rule.CreatedAt) > rule.Expire {
					log.Debugln("Found expired rule, deleting")
					if err := database.Delete(&rule).Error; err != nil {
						log.Warnf("error deleting rule: %s", err)
					}
				}
			}
		}
	}()

	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	app.Post("/rules", func(c *fiber.Ctx) error {
		var r db.Rule
		if err := c.BodyParser(&r); err != nil {
			return response(c, http.StatusUnprocessableEntity, "Invalid request", nil)
		}
		r.ID = ""
		d, err := time.ParseDuration(r.ExpireStr)
		if err != nil {
			return response(c, http.StatusUnprocessableEntity, "Invalid expire duration: "+err.Error(), nil)
		}
		r.Expire = d

		_, err = config.FromJSON(r.Filter)
		if err != nil {
			return response(c, http.StatusUnprocessableEntity, "Invalid filter: "+err.Error(), nil)
		}

		log.Debugf("Creating %+v", r)

		// Add the record
		if err := database.Create(&r).Error; err != nil {
			return internalServerError(c, err)
		}

		return response(c, http.StatusOK, "Rule created", nil)
	})

	app.Get("/rules", func(c *fiber.Ctx) error {
		var rules []db.Rule
		if err := database.Find(&rules).Error; err != nil {
			return internalServerError(c, err)
		}

		log.Debugf("Retreived %+v", rules)

		for i := range rules {
			rules[i].ExpireStr = rules[i].Expire.String()
			log.Debugf("Set expireStr to %s", rules[i].ExpireStr)
		}

		return response(c, http.StatusOK, "Rules retrieved", map[string]interface{}{
			"rules": rules,
		})
	})

	app.Delete("/rules/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		log.Debugf("Deleting rule %s", id)

		if err := database.Delete(&db.Rule{}, "id = ?", id).Error; err != nil {
			return internalServerError(c, err)
		}

		return response(c, http.StatusOK, "Rule deleted", nil)
	})

	app.Get("/generate", func(c *fiber.Ctx) error {
		var rules []db.Rule
		if err := database.Find(&rules).Error; err != nil {
			return internalServerError(c, err)
		}

		var filters []*config.Filter
		for _, rule := range rules {
			filter, err := config.FromJSON(rule.Filter)
			if err != nil {
				// This should never happen because filters are validated before creation
				return internalServerError(c, err)
			}
			filters = append(filters, filter)
		}

		cfg := config.Config{
			Filters: filters,
		}

		return c.Status(http.StatusOK).SendString(cfg.String())
	})

	// go metrics.Listen(metricsListen)

	startupMessage := fmt.Sprintf("Starting API %s on :8080", version)
	sentry.CaptureMessage(startupMessage)
	log.Println(startupMessage)
	log.Fatal(app.Listen(":8080"))
}
