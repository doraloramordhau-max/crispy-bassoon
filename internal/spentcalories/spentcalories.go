package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	data = strings.TrimSpace(data)
	if data == "" {
		return 0, "", 0, errors.New("пустая строка")
	}

	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, errors.New("неверный формат строки данных")
	}

	stepsStr := parts[0]
	actionType := parts[1]
	durationStr := parts[2]

	if strings.Contains(stepsStr, " ") || strings.Contains(durationStr, " ") {
		return 0, "", 0, errors.New("недопустимы пробелы")
	}
	if stepsStr == "+" || stepsStr == "-" {
		return 0, "", 0, errors.New("некорректные шаги")
	}

	// шаги
	var steps int
	var err error
	if strings.HasPrefix(stepsStr, "+") {
		steps, err = strconv.Atoi(stepsStr[1:])
	} else {
		steps, err = strconv.Atoi(stepsStr)
	}
	if err != nil || steps <= 0 {
		return 0, "", 0, errors.New("некорректное количество шагов")
	}

	// попытка стандартного парсинга
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		// поддержка "1.5h" и "30.5m"
		var val float64
		if strings.HasSuffix(durationStr, "h") {
			v := strings.TrimSuffix(durationStr, "h")
			val, err = strconv.ParseFloat(v, 64)
			if err != nil {
				return 0, "", 0, errors.New("ошибка парсинга часов")
			}
			duration = time.Duration(val * float64(time.Hour))
		} else if strings.HasSuffix(durationStr, "m") {
			v := strings.TrimSuffix(durationStr, "m")
			val, err = strconv.ParseFloat(v, 64)
			if err != nil {
				return 0, "", 0, errors.New("ошибка парсинга минут")
			}
			duration = time.Duration(val * float64(time.Minute))
		} else {
			return 0, "", 0, errors.New("некорректная единица измерения")
		}
	}
	if duration <= 0 {
		return 0, "", 0, errors.New("некорректная продолжительность")
	}

	return steps, actionType, duration, nil
}

func distance(steps int, height float64) float64 {
	if steps <= 0 || height <= 0 {
		return 0
	}
	stepLen := height * stepLengthCoefficient
	return (float64(steps) * stepLen) / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if steps <= 0 || height <= 0 || duration <= 0 {
		return 0
	}
	dist := distance(steps, height)
	durationHours := duration.Hours()
	if durationHours <= 0 {
		return 0
	}
	return dist / durationHours
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, action, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)
	durationHours := duration.Hours()

	var calories float64
	switch action {
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	if err != nil {
		return "", err
	}

	result := fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		action, durationHours, dist, speed, calories,
	)

	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные параметры")
	}
	speed := meanSpeed(steps, height, duration)
	durationMin := duration.Minutes()
	calories := (weight * speed * durationMin) / minInH
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные параметры")
	}
	speed := meanSpeed(steps, height, duration)
	durationMin := duration.Minutes()
	calories := ((weight * speed * durationMin) / minInH) * walkingCaloriesCoefficient
	return calories, nil
}
