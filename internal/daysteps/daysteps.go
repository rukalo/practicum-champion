package daysteps

import (
	"fmt"
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
	// разделить строко на слайс строк это делается strings.split()
	splittedData := strings.Split(data, ",")

	if len(splittedData) != 2 {
		return 0, 0, fmt.Errorf("некорректный ввод данных")
	}

	steps, err := strconv.Atoi(splittedData[0])

	if err != nil {
		return 0, 0, err
	}

	if steps <= 0 {

		return 0, 0, fmt.Errorf("число должно быть больше 0")
	}

	walkDuration, err := time.ParseDuration(splittedData[1])
	if err != nil {
		return 0, 0, err
	}

	if walkDuration <= 0 {
		return 0, 0, fmt.Errorf("длительность прогулки не может быть меньше 0")
	}

	return steps, walkDuration, nil

}

func DayActionInfo(data string, weight, height float64) string {
	steps, walkDuration, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	if steps <= 0 {
		return ""
	}

	if weight <= 0 {
		return ""
	}

	if height <= 0 {
		return ""
	}

	distance := float64(steps) * stepLength

	distanceInKm := distance / float64(mInKm)

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, walkDuration)

	return fmt.Sprintf("Количество шагов: %v.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanceInKm, calories)
	// в выводе дописать последний параметр каллорий
}
