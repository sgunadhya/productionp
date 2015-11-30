package productionp

import (
	"log"
	"os"
)

type MPSInput struct {
	forecasts         []int
	minimum_inventory int
	initial_inventory int
	holding_cost      float32
	order_cost        float32
}

type MPSOutput struct {
	plan         []int
	holding_cost float32
	setup_cost   float32
	total_cost   float32
}

func Level(forecasts []int, initial int) int {
	total := 0
	for _, val := range forecasts {
		total += val
	}

	return (total - initial) / len(forecasts)
}

func LevelWithMinimumInventory(forecasts []int, initial int, minimum int) int {
	total := 0
	for _, val := range forecasts {
		total += val
	}
	return (total - (initial - minimum)) / len(forecasts)
}

func Chase(input MPSInput) MPSOutput {
	productions := make([]int, len(input.forecasts))
	for _, val := range input.forecasts {
		if input.initial_inventory <= input.minimum_inventory {
			productions = append(productions, val)
		} else {
			if val < input.initial_inventory {
				input.initial_inventory -= val
				productions = append(productions, 0)
			} else {
				input.initial_inventory = input.minimum_inventory
				productions = append(productions, val-(input.initial_inventory-input.minimum_inventory))
			}
		}
	}
	return MPSOutput{plan: productions}
}

func SilverMeal(mps_input MPSInput) MPSOutput {
	logger := log.New(os.Stderr, "DEBUG: ", log.Ldate|log.Ltime)
	plan := make([]int, len(mps_input.forecasts))
	previous_cost := float32(0.0)
	current_cost := float32(0)
	total_holding_cost := float32(0.0)
	total_setup_cost := float32(0.0)
	trc := float32(0.0)
	order := 0
	logger.Printf("Starting the computation ")

	for i := 0; i < len(mps_input.forecasts); {
		logger.Printf("Computing for period %d current cost : %.2f, previous cost: .2f", i, current_cost, previous_cost)
		trc = mps_input.order_cost
		previous_cost = trc
		total_setup_cost += trc
		order = mps_input.forecasts[i]
		j := i + 1
		unit := 1
		for ; j < len(mps_input.forecasts); j++ {
			logger.Printf("Previous cost : %.2f", previous_cost)
			current_cost = (trc + float32(unit)*mps_input.holding_cost*float32(mps_input.forecasts[j])) / float32(unit+1)
			logger.Printf("Current cost for operation for period %d check %d is : %.2f", i, unit, current_cost)
			logger.Printf("Previous cost for operation for period %d check %d is : %.2f", i, unit, previous_cost)
			if previous_cost < current_cost {
				break
			} else {
				holding_cost := float32(unit) * mps_input.holding_cost * float32(mps_input.forecasts[j])
				trc += holding_cost
				total_holding_cost += holding_cost
				order += mps_input.forecasts[j]
				logger.Printf("Current cost %.2f, holding cost: %.2f , order : %d", trc, holding_cost, order)
			}
			previous_cost = current_cost
			unit++
		}
		plan[i] = order
		i = j
	}
	logger.Printf("%v ", plan)
	logger.Printf("Total holding cost :%.2f, Total Order Cost : %.2f", total_holding_cost, total_setup_cost)
	return MPSOutput{plan: plan, total_cost: total_holding_cost + total_setup_cost,
		holding_cost: total_holding_cost, setup_cost: total_setup_cost}
}

func WagnerWhitin(input MPSInput) MPSOutput {
	dynamic_table := make([][]float32, len(input.forecasts))
	result := make([]int, len(input.forecasts))

	calculate := func(forecasts []int, setup_cost float32, holding_cost float32) float32 {
		total_cost := setup_cost
		for i := 1; i < len(forecasts); i++ {
			total_cost += float32(i) * float32(forecasts[i]) * holding_cost
		}
		return total_cost
	}
	minimum := func(list []float32) float32 {
		min_val := float32(10000000)
		for _, val := range list {
			if val < min_val {
				min_val = val
			}
		}
		return min_val
	}
	minimumWithIndex := func(list []float32) (float32, int) {
		min_val := float32(10000000)
		index := -1
		for i, val := range list {
			if val < min_val {
				min_val = val
				index = i
			}
		}
		return min_val, index
	}

	for i := 0; i < len(input.forecasts); i++ {
		dynamic_table[i] = make([]float32, i+1)
		dynamic_table[i][0] = calculate(input.forecasts[:i+1], input.order_cost, input.holding_cost)
		for j := 0; j < i; j++ {
			dynamic_table[i][j+1] = calculate(input.forecasts[j+1:i+1], input.order_cost, input.holding_cost) + minimum(dynamic_table[j])
		}
	}
	order_quantity := 0
	total_setup_amount := float32(0)

	for j := len(input.forecasts) - 1; j >= 0; j-- {
		_, i := minimumWithIndex(dynamic_table[j])
		order_quantity += input.forecasts[j]
		if i == j {
			result[j] = order_quantity
			order_quantity = 0
			total_setup_amount += input.order_cost
		} else {
			result[j] = 0
		}

	}

	return MPSOutput{plan: result, setup_cost: total_setup_amount}
}
