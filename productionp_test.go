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
	chase := ChaseAlgorithm(MPSInput{Forecasts: forecasts, InitialInventory: initial, MinimumInventory: minimumInventory})
	if 0 != chase.Plan[0] {
		t.Fatalf("Expected %d but got %d", 12, chase.Plan[1])
	}
}

func TestSilverMeal(t *testing.T) {
	mps := MPSInput{Forecasts: []int{120, 130, 100, 150},
		InitialInventory: 0, MinimumInventory: 0, HoldingCost: 5.0, OrderCost: 600.0}

	output := SilverMealAlgorithm(mps)
	if 4 != len(output.Plan) {
		t.Fatalf(" Expected length of plan is : %d but got %d", 4, len(output.Plan))
	}
	if 120 != output.Plan[0] {
		t.Fatalf(" Expected first output is %d but got %d", 120, output.Plan[0])
	}
	if 230 != output.Plan[1] {
		t.Fatalf(" Expected second output is %d but got %d", 230, output.Plan[1])
	}
	if 0 != output.Plan[2] {
		t.Fatalf(" Expected third output is %d but got %d", 0, output.Plan[2])
	}
	if 150 != output.Plan[3] {
		t.Fatalf(" Expected fourth output is %d but got %d", 120, output.Plan[3])
	}

}

func TestWalterWhitin(t *testing.T) {
	mps := MPSInput{Forecasts: []int{120, 130, 100, 150},
		InitialInventory: 0, MinimumInventory: 0, HoldingCost: 5.0, OrderCost: 600.0}
	output := WagnerWhitinAlgorithm(mps)
	if 4 != len(output.Plan) {
		t.Fatalf(" Expected length of plan is : %d but got %d", 4, len(output.Plan))
	}
	if 120 != output.Plan[0] {
		t.Fatalf(" Expected first output is %d but got %d", 120, output.Plan[0])
	}
	if 230 != output.Plan[1] {
		t.Fatalf(" Expected second output is %d but got %d", 230, output.Plan[1])
	}
	if 0 != output.Plan[2] {
		t.Fatalf(" Expected third output is %d but got %d", 0, output.Plan[2])
	}
	if 150 != output.Plan[3] {
		t.Fatalf(" Expected fourth output is %d but got %d", 120, output.Plan[3])
	}
}

func TestAvailableToPromise(t *testing.T) {
	atp := DiscreteAvailableToPromise([]int{1500, 1550, 1600, 1550, 1450, 1400, 1400, 1450, 1400, 1500},
		[]int{6100, 0, 0, 0, 5700, 0, 0, 0, 3000, 0}, []int{800, 750, 1400, 1300, 1000, 970, 980, 850, 750, 700}, 500)

	fmt.Printf("\nThe ATP is : %v \n", atp)
}

func TestProblems(t *testing.T) {
	//mps := MPSInput{Forecasts:[]int{3000,3100,3200,3100,2900,2800,2800,2900,2800,3000}, holding_cost:0.1,order_cost:800}
	mps := MPSInput{Forecasts: []int{300, 280, 200, 190, 160, 210}, HoldingCost: 2, OrderCost: 1000}
	fmt.Printf("%v \n", EOQStrategy(mps, 470))
}
