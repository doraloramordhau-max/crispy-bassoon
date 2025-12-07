package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {

	if strings.TrimSpace(data) == "" {
		return 0, 0, errors.New("пустая строка данных")
	}

	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, errors.New("неверный формат данных — должно быть два параметра")
	}

	stepsStr := parts[0]
	durationStr := parts[1]

	if strings.Contains(stepsStr, " ") || strings.Contains(durationStr, " ") {
		return 0, 0, errors.New("недопустимы пробелы в данных")
	}

	if stepsStr == "+" || stepsStr == "-" {
		return 0, 0, errors.New("некорректное значение шагов")
	}

	steps, err := strconv.Atoi(stepsStr)
	if err != nil || steps <= 0 {
		return 0, 0, errors.New("некорректное количество шагов")
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {

		var val float64
		unit := ""
		if strings.HasSuffix(durationStr, "h") {
			unit = "h"
			valStr := strings.TrimSuffix(durationStr, "h")
			val, err = strconv.ParseFloat(valStr, 64)
			if err != nil {
				return 0, 0, errors.New("некорректный формат часов")
			}
			duration = time.Duration(val * float64(time.Hour))
		} else if strings.HasSuffix(durationStr, "m") {
			unit = "m"
			valStr := strings.TrimSuffix(durationStr, "m")
			val, err = strconv.ParseFloat(valStr, 64)
			if err != nil {
				return 0, 0, errors.New("некорректный формат минут")
			}
			duration = time.Duration(val * float64(time.Minute))
		} else {
			return 0, 0, errors.New("некорректный формат продолжительности")
		}

		if unit != "h" && unit != "m" {
			return 0, 0, errors.New("неверная единица измерения")
		}
	}

	if duration <= 0 {
		return 0, 0, errors.New("продолжительность должна быть больше 0")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	// TODO: реализовать функцию
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	if steps <= 0 {
		return ""
	}

	distanceMeters := float64(steps) * stepLength
	distanceKm := distanceMeters / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println(err)
		return ""
	}

	result := fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps, distanceKm, calories,
	)

	return result
}
