package spentcalories

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep   = 0.65  // средняя длина шага.
	mInKm     = 1000  // количество метров в километре.
	minInH    = 60    // количество минут в часе.
	kmhInMsec = 0.278 // коэффициент для преобразования км/ч в м/с.
	cmInM     = 100   // количество сантиметров в метре.
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parse := strings.Split(data, ",")
	if len(parse) != 3 {
		return 0, "", 0, fmt.Errorf("invalid data format: expected 2 parts, got %d", len(parse))
	}
	steps, err := strconv.Atoi(parse[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("failed to parse steps: %v", err)
	}
	actionsType := parse[1]

	timeAction, err := time.ParseDuration(parse[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("failed to parce timeAction: %v", err)
	}

	return steps, actionsType, timeAction, nil
}

// distance возвращает дистанцию(в километрах), которую преодолел пользователь за время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий (число шагов при ходьбе и беге).
func distance(steps int) float64 {
	return float64(steps) * lenStep / mInKm
}

// meanSpeed возвращает значение средней скорости движения во время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий(число шагов при ходьбе и беге).
// duration time.Duration — длительность тренировки.
func meanSpeed(steps int, duration time.Duration) float64 {
	if duration == 0 {
		return 0
	}
	distKm := distance(steps)
	durationH := duration.Hours()
	return distKm / durationH
}

// ShowTrainingInfo возвращает строку с информацией о тренировке.
//
// Параметры:
//
// data string - строка с данными.
// weight, height float64 — вес и рост пользователя.
func TrainingInfo(data string, weight, height float64) string {

	var Calories, Distance, AvgSpeed float64

	steps, actionsType, timeAction, err := parseTraining(data)
	if err != nil {
		fmt.Errorf("TrainingInfo error parseTraining(data): %w", err)
		return ""
	}
	timeActionHours := timeAction.Hours()

	switch actionsType {
	case "Бег":
		Calories = RunningSpentCalories(steps, weight, timeAction)
		Distance = distance(steps)
		AvgSpeed = meanSpeed(steps, timeAction)

	case "Ходьба":
		Distance = distance(steps)
		AvgSpeed = meanSpeed(steps, timeAction)
		Calories = WalkingSpentCalories(steps, weight, height, timeAction)
	default:
		return "неизвестный тип тренировки"
	}

	return fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		actionsType,     // тип тренировки (string)
		timeActionHours, // длительность в часах (float64)
		Distance,        // дистанция (float64)
		AvgSpeed,        // скорость (float64)
		Calories,        // калории (float64)
	)

}

// Константы для расчета калорий, расходуемых при беге.
const (
	runningCaloriesMeanSpeedMultiplier = 18.0 // множитель средней скорости.
	runningCaloriesMeanSpeedShift      = 20.0 // среднее количество сжигаемых калорий при беге.
)

// RunningSpentCalories возвращает количество потраченных колорий при беге.
//
// Параметры:
//
// steps int - количество шагов.
// weight float64 — вес пользователя.
// duration time.Duration — длительность тренировки.
func RunningSpentCalories(steps int, weight float64, duration time.Duration) float64 {
	return ((runningCaloriesMeanSpeedMultiplier * meanSpeed(steps, duration)) - runningCaloriesMeanSpeedShift) * weight
}

// Константы для расчета калорий, расходуемых при ходьбе.
const (
	walkingCaloriesWeightMultiplier = 0.035 // множитель массы тела.
	walkingSpeedHeightMultiplier    = 0.029 // множитель роста.
)

// WalkingSpentCalories возвращает количество потраченных калорий при ходьбе.
//
// Параметры:
//
// steps int - количество шагов.
// duration time.Duration — длительность тренировки.
// weight float64 — вес пользователя.
// height float64 — рост пользователя.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) float64 {
	hours := duration.Hours()
	return ((walkingCaloriesWeightMultiplier * weight) + (meanSpeed(steps, duration)*meanSpeed(steps, duration)/height)*walkingSpeedHeightMultiplier) * hours * minInH
}
