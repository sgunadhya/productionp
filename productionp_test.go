package productionp

import "testing"

func TestLevel(t *testing.T) {
	forecasts := []int{120, 130, 100, 150}
	initial := 10
	level := Level(forecasts, initial)
	if 122 != level {
		t.Fatalf("Expected %d but got %d", 19, level)
	}

}

func TestLevelWithMinimumInventory(t *testing.T) {
	forecasts := []int{120, 130, 100, 150}
	initial := 10
	minimum_invetory := 10
	level := LevelWithMinimumInventory(forecasts, initial, minimum_invetory)
	if 125 != level {
		t.Fatalf("Expected %d but got %d", 19, level)
	}

}
