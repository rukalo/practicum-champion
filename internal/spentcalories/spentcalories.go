package spentcalories

import (
	"fmt"
	"log"
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
	splittedTrainingData := strings.Split(data, ",")

	if len(splittedTrainingData) == 3 {
		steps, err := strconv.Atoi(splittedTrainingData[0])
		if err != nil {
			return 0, "", 0, err
		}

		walkDuration, err := time.ParseDuration(splittedTrainingData[2])
		if err != nil {
			return 0, "", 0, err
		}

		return steps, splittedTrainingData[1], walkDuration, nil
	}

	return 0, "", 0, nil
}

func distance(steps int, height float64) float64 {
	stepDistance := height * stepLengthCoefficient

	stepsInDistance := float64(steps) * stepDistance

	return stepsInDistance / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	walkDistance := distance(steps, height)

	return walkDistance / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, training, walkDuration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	switch training {
	case "Бег":
		runningDistance := distance(steps, height)
		runningMeanSpeed := meanSpeed(steps, height, walkDuration)
		calories, err := RunningSpentCalories(steps, weight, height, walkDuration)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f", training, walkDuration.Hours(), runningDistance, runningMeanSpeed, calories), nil

	case "Ходьба":
		walkDistance := distance(steps, height)
		walkingMeanSpeed := meanSpeed(steps, height, walkDuration)
		calories, err := WalkingSpentCalories(steps, weight, height, walkDuration)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f", training, walkDuration.Hours(), walkDistance, walkingMeanSpeed, calories), nil

	default:
		return "неизвестный тип тренировки", nil
	}

}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps < 0 {
		return 0, fmt.Errorf("0 шагов")
	}

	if weight <= 0 {
		return 0, fmt.Errorf("вес не может быть ниже 0")
	}

	if height <= 0 {
		return 0, fmt.Errorf("рост не может быть ниже 0")
	}

	if duration < 0 {
		return 0, fmt.Errorf("длительность не может быть ниже 0")
	}

	// в целом в этой проверке нет необходимости т.к первая проверка на duration завершит программу
	if steps > 0 && duration <= 0 {
		return 0, fmt.Errorf("некорректные входные данные - duration")
	}

	runningMeanSpeed := meanSpeed(steps, height, duration)

	walkDurationInMin := duration.Minutes()

	return (weight * runningMeanSpeed * walkDurationInMin) / minInH, nil

}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps < 0 {
		return 0, fmt.Errorf("шагов не может быть ниже 0")
	}

	if weight <= 0 {
		return 0, fmt.Errorf("вес не может быть ниже 0")
	}

	if height <= 0 {
		return 0, fmt.Errorf("рост не может быть ниже 0")
	}

	if duration < 0 {
		return 0, fmt.Errorf("длительность не может быть ниже 0")
	}

	walkingMeanSpeed := meanSpeed(steps, height, duration)

	walkDurationInMin := duration.Minutes()

	resultCalories := (weight * walkingMeanSpeed * walkDurationInMin) / minInH

	return resultCalories * walkingCaloriesCoefficient, nil
}
