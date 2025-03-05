package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	calc "github.com/Evrard-ro/track-go/internal/spentcalories"
)

var (
	StepLength = 0.65 // длина шага в метрах
)

func parsePackage(data string) (int, time.Duration, error) {
	pars := strings.Split(data, ",")
	if len(pars) != 2 {
		return 0, 0, fmt.Errorf("invalid data format: expected 2 parts, got %d", len(pars))
	}
	steps, err := strconv.Atoi(pars[0])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse steps: %v", err)
	}
	duration, err := time.ParseDuration(pars[1])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse duration: %v", err)
	}

	return steps, duration, nil
}

// DayActionInfo обрабатывает входящий пакет, который передаётся в
// виде строки в параметре data. Параметр storage содержит пакеты за текущий день.
// Если время пакета относится к новым суткам, storage предварительно
// очищается.
// Если пакет валидный, он добавляется в слайс storage, который возвращает
// функция. Если пакет невалидный, storage возвращается без изменений.
func DayActionInfo(data string, weight, height float64) string {
	steps, timeTrain, err := parsePackage(data)
	if err != nil {
		fmt.Errorf("data format error: %v", err)
		return ""
	}
	if steps <= 0 {
		fmt.Errorf("number of steps is not positive")
		return ""
	}
	distance := float64(steps) * StepLength / 1000
	calories := calc.WalkingSpentCalories(steps, weight, height, timeTrain)
	return fmt.Sprintf("Колличество шагов: %d\n Дистанция составила %.2fкм.\n Вы сожгли %.2f ккал.", steps, distance, calories)
}
