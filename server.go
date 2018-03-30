package funding

type FundServer struct {
	// Commands chan interface{}
	// commands chan interface{}
	// fund     Fund

	commands chan TransactionCommand
	fund     *Fund
}

type WithdrawCommand struct {
	Amount int
}

type BalanceCommand struct {
	Response chan int
}

// Typedef the callback for readability
// type Transactor func(fund *Fund)
type Transactor func(interface{})

// Add a new command type with a callback and a semaphore channel
type TransactionCommand struct {
	Transactor Transactor
	Done       chan bool
}

func NewFundServer(initialBalance int) *FundServer {
	server := &FundServer{
		// make() creates builtins like channels, maps, and slices
		// Commands: make(chan interface{}),
		fund: NewFund(initialBalance),
	}

	// Spawn off the server's main loop immediately
	go server.loop()
	return server
}

func (s *FundServer) loop() {
	// for command := range s.commands {

	// // command is just an interface{}, but we can check its real type
	// switch command.(type) {
	//
	// case WithdrawCommand:
	// 	// And then use a "type assertion" to convert it
	// 	withdrawal := command.(WithdrawCommand)
	// 	s.fund.Withdraw(withdrawal.Amount)
	//
	// case BalanceCommand:
	// 	getBalance := command.(BalanceCommand)
	// 	balance := s.fund.Balance()
	// 	getBalance.Response <- balance
	//
	// case TransactionCommand:
	// 	transaction := command.(TransactionCommand)
	// 	transaction.Transactor(s.fund)
	// 	transaction.Done <- true
	//
	// default:
	// 	panic(fmt.Sprintf("Unrecognized command: %v", command))
	// }
	// }
	for transaction := range s.commands {
		// Now we don't need any type-switch mess
		transaction.Transactor(s.fund)
		transaction.Done <- true
	}

}

// func (s *FundServer) Balance() int {
// 	responseChan := make(chan int)
// 	s.commands <- BalanceCommand{Response: responseChan}
// 	return <-responseChan
// }

func (s *FundServer) Balance() int {
	var balance int
	s.Transact(func(f *Fund) {
		balance = f.Balance()
	})
	return balance
}

// func (s *FundServer) Withdraw(amount int) {
// 	s.commands <- WithdrawCommand{Amount: amount}
// }

func (s *FundServer) Withdraw(amount int) {
	s.Transact(func(f *Fund) {
		f.Withdraw(amount)
	})
}

// Wrap it up neatly in an API method, like the other commands
func (s *FundServer) Transact(transactor Transactor) {
	command := TransactionCommand{
		Transactor: transactor,
		Done:       make(chan bool),
	}
	s.commands <- command
	<-command.Done
}
