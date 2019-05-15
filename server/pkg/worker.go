package server

type AddItemAction struct {
	Code         Code
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
	}

	return worker
}

func (w *Worker) Start() {
	go w.Run()
}

func (w *Worker) GetId() string {
	return string(w.Basket.UUID.Bytes())
}

func (w *Worker) Run() {
	for {
		select {
		case addItemAction := <-w.AddItem:
			code := addItemAction.Code
			item, ok := Items[code]
			if !ok {
				addItemAction.ErrorChan <- &UnkownItemError{}
				return
			}

			w.Basket.AddItem(item)
			addItemAction.GetTotalChan <- w.Basket.GetTotal()

		case closeAction := <-w.Close:
			closeAction.GetTotalChan <- w.Basket.GetTotal()
			return
		}
	}
}
