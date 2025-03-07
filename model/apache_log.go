package model

import (
	"reflect"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type ApacheLog struct {
	gorm.Model
	ClientIP        string           `gorm:"size:45" log:"a" json:"client_ip"`         // %a
	LocalIP         string           `gorm:"size:45" log:"A" json:"local_ip"`          // %A
	VirtualHostName string           `gorm:"size:255" log:"v" json:"virtualhost_name"` // %v
	Time            time.Time        `gorm:"autoCreateTime" log:"t" json:"time"`       // %{cu}t
	Headers         ApacheLogHeaders `gorm:"serializer:json"`
}

func LoopThroughReflection(entryValue reflect.Value, reflect_key string, key string, value string) {
	entryType := entryValue.Type()
	for j := 0; j < entryType.NumField(); j++ {
		structField := entryType.Field(j)
		tag := structField.Tag.Get(reflect_key)

		if tag == key {
			fieldValue := entryValue.Field(j)

			// Convert field value to appropriate type
			switch fieldValue.Kind() {
			case reflect.Int:
				intVal, _ := strconv.Atoi(value)
				fieldValue.SetInt(int64(intVal))
			case reflect.String:
				fieldValue.SetString(value)
			case reflect.Struct:
				if structField.Type == reflect.TypeOf(time.Time{}) {
					parsedTime, _ := time.Parse("2006-01-02 15:04:05.000000 -0700", value)
					fieldValue.Set(reflect.ValueOf(parsedTime))
				}
			}
		}

		if structField.Anonymous {
			embeddedValue := entryValue.Field(j)

			LoopThroughReflection(embeddedValue, reflect_key, key, value)
		}
	}
}
