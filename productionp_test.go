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

func TestChase(t *testing.T) {
	forecasts := []int{120, 130, 100, 150}
	initial := 10
	minimum_invetory := 10
	chase := Chase(forecasts, initial, minimum_invetory)
	if 12 != chase[0] {
		t.Fatalf("Expected %d but got %d", 12, chase[1])
	}
}

func TestSilverMeal(t *testing.T) {
	mps := MPSInput{forecasts: []int{120, 130, 100, 150},
		initial_inventory: 0, minimum_inventory: 0, holding_cost: 5.0, order_cost: 600.0}

	output := SilverMeal(mps)
	if 4 != len(output.plan) {
		t.Fatalf(" Expected length of plan is : %d but got %d", 4, len(output.plan))
	}
	if 120 != output.plan[0] {
		t.Fatalf(" Expected first output is %d but got %d", 120, output.plan[0])
	}
	if 230 != output.plan[1] {
		t.Fatalf(" Expected second output is %d but got %d", 230, output.plan[1])
	}
	if 0 != output.plan[2] {
		t.Fatalf(" Expected third output is %d but got %d", 0, output.plan[2])
	}
	if 150 != output.plan[3] {
		t.Fatalf(" Expected fourth output is %d but got %d", 120, output.plan[3])
	}

}
