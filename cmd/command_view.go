package cmd

import (
	"fmt"

	"github.com/charmbracelet/glamour"
	"github.com/harljos/gymplanr/internal/database"
)

func viewCmd(cfg *config, user database.User) error {
	databaseDays, err := cfg.getDaysByUser(user)
	if err != nil {
		return err
	}
	if databaseDays == nil {
		fmt.Println("Workout plan has not been found use the 'generate' command to get one")
		return nil
	}

	var in string

	for _, day := range databaseDays {
		in = fmt.Sprintf(`%s
		
# %s
| Exercise | Sets | Reps |
|----------|------|------|`, in, day.Name)
		databaseExercises, err := cfg.getExercisesByDay(day)
		if err != nil {
			return err
		}
		for _, exercise := range databaseExercises {
			in = fmt.Sprintf(`%s
| %v | %v | %v |`, in, exercise.Name, exercise.Sets.Int32, exercise.Repetitions.Int32)
		}
	}

	out, err := glamour.Render(in, "dark")
	if err != nil {
		return err
	}

	fmt.Println(out)
	return nil
}
