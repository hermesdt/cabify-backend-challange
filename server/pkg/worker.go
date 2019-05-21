package server

type AddItemAction struct {
	Item         Item
	GetTotalChan chan<- Total
	ErrorChan    chan<- error
}

type CloseAction struct {
	GetTotalChan chan<- Total
}

type Worker struct {
	Basket  *Basket
	AddItem chan AddItemAction
	Close   chan CloseAction
}

func NewWorker() *Worker {
	worker := &Worker{
		Basket:  NewBasket(),
		AddItem: make(chan AddItemAction),
		Close:   make(chan CloseAction),
	}

	return worker
}

func (w *Worker) Start() {
	go w.Run()
}

func (w *Worker) GetId() string {
	return w.Basket.UUID.String()
}

func (w *Worker) Run() {
	for {
		select {
		case addItemAction := <-w.AddItem:
			w.Basket.AddItem(addItemAction.Item)
			addItemAction.GetTotalChan <- w.Basket.GetTotal()

		case closeAction := <-w.Close:
			closeAction.GetTotalChan <- w.Basket.GetTotal()
			return
		}
	}
}
