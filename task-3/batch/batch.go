package batch

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type reqData struct {
	data string
	id   int
}

type Service struct {
	ip             string
	port           string
	size           int
	timeout        time.Duration
	ticker         *time.Ticker
	m              *sync.Mutex
	lockMap        map[int]chan bool
	chData         chan reqData
	currentRequest int
	db             ExternalDB
}

func NewService(db ExternalDB, ip, port string, size int, timeout time.Duration) *Service {
	if size < 1 {
		panic("bad size data")
	}
	if timeout.Seconds() < 1 {
		panic("bad timeout")
	}
	b := &Service{
		ip:             ip,
		port:           port,
		size:           size,
		timeout:        timeout,
		ticker:         time.NewTicker(timeout),
		m:              new(sync.Mutex),
		chData:         make(chan reqData, size),
		currentRequest: 0,
		lockMap:        make(map[int]chan bool),
		db:             db,
	}

	return b
}

func (b *Service) Start() {
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", b.httpHandler)
		err := http.ListenAndServe(fmt.Sprintf("%s:%s", b.ip, b.port), mux)
		if err != nil {
			panic(err)
		}
	}()
	go b.loop()
}

func (b *Service) flush(elems []reqData, count int) {
	b.m.Lock()
	defer b.m.Unlock()
	dataArr := make([]string, count)
	for i := 0; i < count; i++ {
		dataArr[i] = elems[i].data
	}

	b.db.SaveBatch(dataArr, count)

	for i := 0; i < count; i++ {
		e := elems[i]
		//fmt.Println(e.data)
		b.lockMap[e.id] <- true
	}
}

func (b *Service) loop() {

	arr := make([]reqData, b.size)
	elem := 0

	for {
		select {
		case <-b.ticker.C:
			tmp := make([]reqData, b.size)
			copy(tmp, arr)
			b.flush(tmp, elem)
			elem = 0
		case data := <-b.chData:
			arr[elem] = data
			if elem+1 == b.size {
				tmp := make([]reqData, b.size)
				copy(tmp, arr)
				b.flush(tmp, elem+1)
				elem = 0
			} else {
				elem++
			}
		}
	}
}

func (b *Service) Push(data string) error {
	b.m.Lock()
	reqId := b.currentRequest
	b.currentRequest++
	b.lockMap[reqId] = make(chan bool)
	b.m.Unlock()

	b.chData <- reqData{
		data: data,
		id:   reqId,
	}

	<-b.lockMap[reqId]

	b.m.Lock()
	delete(b.lockMap, reqId)
	b.m.Unlock()

	return nil
}
