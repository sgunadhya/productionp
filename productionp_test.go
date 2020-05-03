package productionp

import (
	"fmt"
	"testing"
)

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
	minimumInventory := 10
	level := LevelWithMinimumInventory(forecasts, initial, minimumInventory)
	if 125 != level {
		t.Fatalf("Expected %d but got %d", 19, level)
	}

}

func TestChase(t *testing.T) {
	forecasts := []int{120, 130, 100, 150}
	initial := 10
	minimumInventory := 10
	chase := ChaseAlgorithm(MPSInput{forecasts: forecasts, initialInventory: initial, minimumInventory: minimumInventory})
	if 0 != chase.plan[0] {
		t.Fatalf("Expected %d but got %d", 12, chase.plan[1])
	}
}

func TestSilverMeal(t *testing.T) {
	mps := MPSInput{forecasts: []int{120, 130, 100, 150},
		initialInventory: 0, minimumInventory: 0, holdingCost: 5.0, orderCost: 600.0}

	output := SilverMealAlgorithm(mps)
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

func TestWalterWhitin(t *testing.T) {
	mps := MPSInput{forecasts: []int{120, 130, 100, 150},
		initialInventory: 0, minimumInventory: 0, holdingCost: 5.0, orderCost: 600.0}
	output := WagnerWhitinAlgorithm(mps)
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

func TestAvailableToPromise(t *testing.T) {
	atp := DiscreteAvailableToPromise([]int{1500, 1550, 1600, 1550, 1450, 1400, 1400, 1450, 1400, 1500},
		[]int{6100, 0, 0, 0, 5700, 0, 0, 0, 3000, 0}, []int{800, 750, 1400, 1300, 1000, 970, 980, 850, 750, 700}, 500)

	fmt.Printf("\nThe ATP is : %v \n", atp)
}

func TestProblems(t *testing.T) {
	//mps := MPSInput{forecasts:[]int{3000,3100,3200,3100,2900,2800,2800,2900,2800,3000}, holding_cost:0.1,order_cost:800}
	mps := MPSInput{forecasts: []int{300, 280, 200, 190, 160, 210}, holdingCost: 2, orderCost: 1000}
	fmt.Printf("%v \n", EOQStrategy(mps, 470))
}
