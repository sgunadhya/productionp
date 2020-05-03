package productionp

import (
	"fmt"
	"log"
	"os"
)

type MPSInput struct {
	Forecasts        []int
	MinimumInventory int
	InitialInventory int
	HoldingCost      float32
	OrderCost        float32
}

type MPSOutput struct {
	Plan        []int
	HoldingCost float32
	SetupCost   float32
	TotalCost   float32
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

func ChaseAlgorithm(input MPSInput) MPSOutput {
	productionPlans := make([]int, len(input.Forecasts))
	orderCost := float32(0.0)
	holdingCost := float32(0.0)
	holdingCounter := 1
	for _, val := range input.Forecasts {
		if input.InitialInventory <= input.MinimumInventory {
			productionPlans = append(productionPlans, val)
			orderCost += input.OrderCost
			holdingCounter = 1
		} else {
			if val < input.InitialInventory {
				input.InitialInventory -= val
				productionPlans = append(productionPlans, 0)
				holdingCounter++
				holdingCost += float32(holdingCounter) * input.HoldingCost * float32(val)
			} else {
				input.InitialInventory = input.MinimumInventory
				productionPlans = append(productionPlans, val-(input.InitialInventory-input.MinimumInventory))
				orderCost += input.OrderCost
				holdingCounter = 1
			}
		}
	}
	return MPSOutput{Plan: productionPlans, HoldingCost: holdingCost, SetupCost: orderCost, TotalCost: holdingCost + orderCost}
}

func EOQStrategy(input MPSInput, eoq int) MPSOutput {
	productionPlans := make([]int, len(input.Forecasts))
	inventoryOnHand := input.InitialInventory
	holdingCounter := 0
	totalHoldingCost := float32(0.0)
	totalSetupCost := float32(0.0)

	for i := 0; i < len(input.Forecasts); i++ {
		if inventoryOnHand < input.Forecasts[i] {
			productionPlans[i] = eoq
			totalSetupCost += input.OrderCost
			holdingCounter = 1
			inventoryOnHand += eoq
		} else {
			productionPlans[i] = 0
			holdingCounter++
			totalHoldingCost += float32(holdingCounter) * input.HoldingCost * float32(input.Forecasts[i])
		}
		inventoryOnHand -= input.Forecasts[i]

	}

	fmt.Printf("%v   \n", productionPlans)
	totalCost := totalHoldingCost + totalSetupCost
	return MPSOutput{Plan: productionPlans, TotalCost: totalCost}
}

func SilverMealAlgorithm(mpsInput MPSInput) MPSOutput {
	logger := log.New(os.Stderr, "DEBUG: ", log.Ldate|log.Ltime)
	plan := make([]int, len(mpsInput.Forecasts))
	previousCost := float32(0.0)
	currentCost := float32(0)
	totalHoldingCost := float32(0.0)
	totalSetupCost := float32(0.0)
	trc := float32(0.0)
	order := 0
	logger.Printf("Starting the computation ")

	for i := 0; i < len(mpsInput.Forecasts); {
		logger.Printf("Computing for period %d current cost : %.2f, previous cost: %.2f", i, currentCost, previousCost)
		trc = mpsInput.OrderCost
		previousCost = trc
		totalSetupCost += trc
		order = mpsInput.Forecasts[i]
		j := i + 1
		unit := 1
		for ; j < len(mpsInput.Forecasts); j++ {
			logger.Printf("Previous cost : %.2f", previousCost)
			currentCost = (trc + float32(unit)*mpsInput.HoldingCost*float32(mpsInput.Forecasts[j])) / float32(unit+1)
			logger.Printf("Current cost for operation for period %d check %d is : %.2f", i, unit, currentCost)
			logger.Printf("Previous cost for operation for period %d check %d is : %.2f", i, unit, previousCost)
			if previousCost < currentCost {
				break
			} else {
				holdingCost := float32(unit) * mpsInput.HoldingCost * float32(mpsInput.Forecasts[j])
				trc += holdingCost
				totalHoldingCost += holdingCost
				order += mpsInput.Forecasts[j]
				logger.Printf("Current cost %.2f, holding cost: %.2f , order : %d", trc, holdingCost, order)
			}
			previousCost = currentCost
			unit++
		}
		plan[i] = order
		i = j
	}
	logger.Printf("%v ", plan)
	logger.Printf("Total holding cost :%.2f, Total Order Cost : %.2f", totalHoldingCost, totalSetupCost)
	return MPSOutput{Plan: plan, TotalCost: totalHoldingCost + totalSetupCost,
		HoldingCost: totalHoldingCost, SetupCost: totalSetupCost}
}

func WagnerWhitinAlgorithm(input MPSInput) MPSOutput {
	dynamicTable := make([][]float32, len(input.Forecasts))
	productionPlans := make([]int, len(input.Forecasts))

	calculate := func(forecasts []int, setupCost float32, holdingCost float32) float32 {
		totalCost := setupCost
		for i := 1; i < len(forecasts); i++ {
			totalCost += float32(i) * float32(forecasts[i]) * holdingCost
		}
		return totalCost
	}
	minimumWithIndex := func(list []float32) (float32, int) {
		minVal := float32(10000000)
		index := -1
		for i, val := range list {
			if val < minVal {
				minVal = val
				index = i
			}
		}
		return minVal, index
	}

	for i := 0; i < len(input.Forecasts); i++ {
		dynamicTable[i] = make([]float32, i+1)
		dynamicTable[i][0] = calculate(input.Forecasts[:i+1], input.OrderCost, input.HoldingCost)
		for j := 0; j < i; j++ {
			minimumVal, _ := minimumWithIndex(dynamicTable[j])
			dynamicTable[i][j+1] = calculate(input.Forecasts[j+1:i+1], input.OrderCost, input.HoldingCost) + minimumVal
		}
	}
	orderQuantity := 0
	totalSetupAmount := float32(0)
	fmt.Printf("%v \n", dynamicTable)
	minIndex := len(input.Forecasts)
	totalCost, _ := minimumWithIndex(dynamicTable[len(input.Forecasts)-1])

	for j := len(input.Forecasts) - 1; j >= 0; j-- {
		if minIndex < len(dynamicTable[j])-1 {
			productionPlans[j] = 0
			continue
		} else {
			_, i := minimumWithIndex(dynamicTable[j])
			orderQuantity += input.Forecasts[j]
			minIndex = i
			productionPlans[j] = orderQuantity
			orderQuantity = 0
			totalSetupAmount += input.OrderCost
		}
	}

	return MPSOutput{Plan: productionPlans, SetupCost: totalSetupAmount, TotalCost: totalCost}
}

func DiscreteAvailableToPromise(forecasts []int, productionPlans []int, committedOrders []int, inventoryOnHand int) []int {
	atp := make([]int, len(forecasts))
	atp[0] = inventoryOnHand
	committedNextRun := 0

	for i := 0; i < len(forecasts); {
		committedNextRun = committedOrders[i]
		j := i + 1
		for ; j < len(forecasts); j++ {
			fmt.Printf("j ....%d", j)
			if productionPlans[j] != 0 {
				break
			} else {
				committedNextRun += committedOrders[j]
			}
		}
		atp[i] += productionPlans[i] - committedNextRun
		i = j
		committedNextRun = 0
	}

	return atp

}
