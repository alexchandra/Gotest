package funding

import (
	"sync"
	"testing"
)

const WORKERS = 10

func BenchmarkWithdrawals(b *testing.B) {
	// Skip N = 1
	if b.N < WORKERS {
		return
	}

	// Add as many dollars as we have iterations this run
	//fund := NewFund(b.N)

	// Add new server
	server := NewFundServer(b.N)

	// Casually assume b.N divides cleanly
	dollarsPerFounder := b.N / WORKERS

	// WaitGroup structs don't need to be initialized
	// (their "zero value" is ready to use).
	// So, we just declare one and then use it.
	var wg sync.WaitGroup

	for i := 0; i < WORKERS; i++ {
		// Let the waitgroup know we're adding a goroutine
		wg.Add(1)

		// Spawn off a founder worker, as a closure
		go func() {
			// Mark this worker done when the function finishes
			defer wg.Done()

			pizzaTime := false

			for i := 0; i < dollarsPerFounder; i++ {
				// Stop when we're down to pizza money
				// if server.Balance() <= 10 {
				// 	break
				// }

				// server.Commands <- WithdrawCommand{Amount: 1}
				// server.Withdraw(1)

				// server.Transact(func(fund *Fund) {
				// 	if fund.Balance() <= 10 {
				// 		// Set it in the outside scope
				// 		pizzaTime = true
				// 		return
				// 	}
				// 	fund.Withdraw(1)
				// })

				server.Transact(func(managedValue interface{}) {
					fund := managedValue.(*Fund)
					// Do stuff with fund ...
					fund.Withdraw(1)
				})

				if pizzaTime {
					break
				}
			}

		}() // Remember to call the closure!
	}

	// Wait for all the workers to finish
	wg.Wait()

	// balanceResponseChan := make(chan int)
	// server.Commands <- BalanceCommand{Response: balanceResponseChan}
	// balance := <-balanceResponseChan
	balance := server.Balance()

	if balance != 10 {
		b.Error("Balance wasn't ten dollars:", balance)
	}
}
