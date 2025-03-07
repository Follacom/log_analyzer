package lib

import (
	"bufio"
	"fmt"
	"log_analyzer/model"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var MuError sync.Mutex
var MuDatabase sync.Mutex

func LogError(primaryErr error) {
	MuError.Lock()
	errorFile, err := os.OpenFile("./log_analyzer.error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer errorFile.Close()

	_, err = errorFile.Write([]byte(primaryErr.Error() + "\r"))
	if err != nil {
		return
	}
	MuError.Unlock()
}

func ListenAccess() {
	// Loop all access log files
	accessPathFiles := viper.GetViper().GetStringSlice("scan.access.path")
	for _, accessPathFile := range accessPathFiles {
		if _, err := os.Stat(accessPathFile); err != nil {
			LogError(err)
			return
		}
		go listenAccessFile(accessPathFile)
	}
}

func listenAccessFile(logFile string) {
	// Configure the batch size for import
	batchSize := viper.GetViper().GetInt("database.batch_size")

	for {
		parsedLog := ParseFile(logFile)
		if parsedLog.Len() != 0 {
			scanner := bufio.NewScanner(&parsedLog)

			db, err := RotateDatabase()
			if err != nil {
				LogError(err)
				return
			}
			defer func() {
				sql, err := db.DB()
				if err != nil {
					LogError(err)
				}
				sql.Close()
			}()
			// Migrate the schema
			db.AutoMigrate(&model.ApacheAccessLog{})
			db.Transaction(func(tx *gorm.DB) error {
				logs := make([]model.ApacheAccessLog, 0)
				for scanner.Scan() {
					line := scanner.Text()

					if viper.GetViper().GetBool("scan.access.keep_logs") {
						dir := filepath.Dir(viper.GetViper().GetString("database.url"))
						tmp, err := os.OpenFile(RotateFile(filepath.Join(dir, "log_analyzer_error.log")), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
						if err != nil {
							LogError(err)
							return err
						}
						defer tmp.Close()

						_, err = tmp.Write([]byte(line + "\r"))
						if err != nil {
							LogError(err)
							return err
						}

						tmp.Close()
					}

					accessLog := new(model.ApacheAccessLog)

					if err := accessLog.Parse(line); err == nil {
						// db.Create(accessLog)
						logs = append(logs, *accessLog)

						if len(logs) >= batchSize {
							if err := tx.CreateInBatches(logs, batchSize).Error; err != nil {
								LogError(err)
								return err
							}
							logs = nil // Reset slice after batch insert
						}
					}
				}

				// Insert any remaining logs
				if len(logs) > 0 {
					if err := tx.CreateInBatches(logs, batchSize).Error; err != nil {
						LogError(err)
						return err
					}
				}

				if err := scanner.Err(); err != nil {
					LogError(err)
					return err
				}

				return nil
			})

			sql, _ := db.DB()
			sql.Close()
		}

		time.Sleep(viper.GetViper().GetDuration("scan.interval"))
	}
}

func ListenError() {
	// Loop all access log files
	errorPathFiles := viper.GetViper().GetStringSlice("scan.error.path")
	for _, errorPathFile := range errorPathFiles {
		_, err := os.Stat(errorPathFile)
		if err != nil && err != os.ErrNotExist {
			LogError(err)
		} else {
			go listenErrorFile(errorPathFile)
		}
	}
}

func listenErrorFile(logFile string) {
	// Configure the batch size for import
	batchSize := viper.GetViper().GetInt("database.batch_size")
	for {
		parsedLog := ParseFile(logFile)
		if parsedLog.Len() != 0 {
			scanner := bufio.NewScanner(&parsedLog)

			db, err := RotateDatabase()
			if err != nil {
				LogError(err)
				return
			}
			defer func() {
				sql, _ := db.DB()
				sql.Close()
			}()
			// Migrate the schema
			db.AutoMigrate(&model.ApacheErrorLog{})
			db.Transaction(func(tx *gorm.DB) error {
				logs := make([]model.ApacheErrorLog, 0)
				for scanner.Scan() {
					line := scanner.Text()

					if viper.GetViper().GetBool("scan.error.keep_logs") {
						dir := filepath.Dir(viper.GetViper().GetString("database.url"))
						tmp, err := os.OpenFile(RotateFile(filepath.Join(dir, "log_analyzer_error.log")), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
						if err != nil {
							LogError(err)
							return err
						}
						defer tmp.Close()

						_, err = tmp.Write([]byte(line + "\r"))
						if err != nil {
							LogError(err)
							return err
						}

						tmp.Close()
					}

					errorLog := new(model.ApacheErrorLog)

					switch line {
					case "The 'Apache 2.4' service is restarting.":
						line = fmt.Sprintf("[l:\"notice\"] [m:\"httpd_service\"] [M:\"%s\"] [v:\"localhost\"]", line)
					case "Starting the 'Apache 2.4' service":
						line = fmt.Sprintf("[l:\"notice\"] [m:\"httpd_service\"] [M:\"%s\"] [v:\"localhost\"]", line)
					case "The 'Apache 2.4' service is running.":
						line = fmt.Sprintf("[l:\"notice\"] [m:\"httpd_service\"] [M:\"%s\"] [v:\"localhost\"]", line)
					default:
					}
					if err := errorLog.Parse(line); err == nil {
						// db.Create(accessLog)
						logs = append(logs, *errorLog)

						if len(logs) >= batchSize {
							if err := tx.CreateInBatches(logs, batchSize).Error; err != nil {
								LogError(err)
								return err
							}
							logs = nil // Reset slice after batch insert
						}
					}
				}

				// Insert any remaining logs
				if len(logs) > 0 {
					if err := tx.CreateInBatches(logs, batchSize).Error; err != nil {
						LogError(err)
						return err
					}
				}

				if err := scanner.Err(); err != nil {
					LogError(err)
					return err
				}

				return nil
			})

			sql, _ := db.DB()
			sql.Close()
		}

		time.Sleep(viper.GetViper().GetDuration("scan.interval"))
	}
}
