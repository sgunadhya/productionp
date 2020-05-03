package productionp

import (
	"fmt"
	"log"
	"os"
)

type MPSInput struct {
	forecasts        []int
	minimumInventory int
	initialInventory int
	holdingCost      float32
	orderCost        float32
}

type MPSOutput struct {
	plan        []int
	holdingCost float32
	setupCost   float32
	totalCost   float32
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
	productionPlans := make([]int, len(input.forecasts))
	orderCost := float32(0.0)
	holdingCost := float32(0.0)
	holdingCounter := 1
	for _, val := range input.forecasts {
		if input.initialInventory <= input.minimumInventory {
			productionPlans = append(productionPlans, val)
			orderCost += input.orderCost
			holdingCounter = 1
		} else {
			if val < input.initialInventory {
				input.initialInventory -= val
				productionPlans = append(productionPlans, 0)
				holdingCounter++
				holdingCost += float32(holdingCounter) * input.holdingCost * float32(val)
			} else {
				input.initialInventory = input.minimumInventory
				productionPlans = append(productionPlans, val-(input.initialInventory-input.minimumInventory))
				orderCost += input.orderCost
				holdingCounter = 1
			}
		}
	}
	return MPSOutput{plan: productionPlans, holdingCost: holdingCost, setupCost: orderCost, totalCost: holdingCost + orderCost}
}

func EOQStrategy(input MPSInput, eoq int) MPSOutput {
	productionPlans := make([]int, len(input.forecasts))
	inventoryOnHand := input.initialInventory
	holdingCounter := 0
	totalHoldingCost := float32(0.0)
	totalSetupCost := float32(0.0)

	for i := 0; i < len(input.forecasts); i++ {
		if inventoryOnHand < input.forecasts[i] {
			productionPlans[i] = eoq
			totalSetupCost += input.orderCost
			holdingCounter = 1
			inventoryOnHand += eoq
		} else {
			productionPlans[i] = 0
			holdingCounter++
			totalHoldingCost += float32(holdingCounter) * input.holdingCost * float32(input.forecasts[i])
		}
		inventoryOnHand -= input.forecasts[i]

	}

	fmt.Printf("%v   \n", productionPlans)
	totalCost := totalHoldingCost + totalSetupCost
	return MPSOutput{plan: productionPlans, totalCost: totalCost}
}

func SilverMealAlgorithm(mpsInput MPSInput) MPSOutput {
	logger := log.New(os.Stderr, "DEBUG: ", log.Ldate|log.Ltime)
	plan := make([]int, len(mpsInput.forecasts))
	previousCost := float32(0.0)
	currentCost := float32(0)
	totalHoldingCost := float32(0.0)
	totalSetupCost := float32(0.0)
	trc := float32(0.0)
	order := 0
	logger.Printf("Starting the computation ")

	for i := 0; i < len(mpsInput.forecasts); {
		logger.Printf("Computing for period %d current cost : %.2f, previous cost: %.2f", i, currentCost, previousCost)
		trc = mpsInput.orderCost
		previousCost = trc
		totalSetupCost += trc
		order = mpsInput.forecasts[i]
		j := i + 1
		unit := 1
		for ; j < len(mpsInput.forecasts); j++ {
			logger.Printf("Previous cost : %.2f", previousCost)
			currentCost = (trc + float32(unit)*mpsInput.holdingCost*float32(mpsInput.forecasts[j])) / float32(unit+1)
			logger.Printf("Current cost for operation for period %d check %d is : %.2f", i, unit, currentCost)
			logger.Printf("Previous cost for operation for period %d check %d is : %.2f", i, unit, previousCost)
			if previousCost < currentCost {
				break
			} else {
				holdingCost := float32(unit) * mpsInput.holdingCost * float32(mpsInput.forecasts[j])
				trc += holdingCost
				totalHoldingCost += holdingCost
				order += mpsInput.forecasts[j]
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
	return MPSOutput{plan: plan, totalCost: totalHoldingCost + totalSetupCost,
		holdingCost: totalHoldingCost, setupCost: totalSetupCost}
}

func WagnerWhitinAlgorithm(input MPSInput) MPSOutput {
	dynamicTable := make([][]float32, len(input.forecasts))
	productionPlans := make([]int, len(input.forecasts))

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

	for i := 0; i < len(input.forecasts); i++ {
		dynamicTable[i] = make([]float32, i+1)
		dynamicTable[i][0] = calculate(input.forecasts[:i+1], input.orderCost, input.holdingCost)
		for j := 0; j < i; j++ {
			minimumVal, _ := minimumWithIndex(dynamicTable[j])
			dynamicTable[i][j+1] = calculate(input.forecasts[j+1:i+1], input.orderCost, input.holdingCost) + minimumVal
		}
	}
	orderQuantity := 0
	totalSetupAmount := float32(0)
	fmt.Printf("%v \n", dynamicTable)
	minIndex := len(input.forecasts)
	totalCost, _ := minimumWithIndex(dynamicTable[len(input.forecasts)-1])

	for j := len(input.forecasts) - 1; j >= 0; j-- {
		if minIndex < len(dynamicTable[j])-1 {
			productionPlans[j] = 0
			continue
		} else {
			_, i := minimumWithIndex(dynamicTable[j])
			orderQuantity += input.forecasts[j]
			minIndex = i
			productionPlans[j] = orderQuantity
			orderQuantity = 0
			totalSetupAmount += input.orderCost
		}
	}

	return MPSOutput{plan: productionPlans, setupCost: totalSetupAmount, totalCost: totalCost}
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
