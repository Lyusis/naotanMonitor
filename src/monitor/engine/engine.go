package engine

import (
	"api"
	"config"
	"monitor/model"
)

// Run 多线程调度引擎/**
func (engine *ConcurrentEngine) Run(seeds ...Request) {

	out := make(chan Result)
	engine.Scheduler.Run()

	for i := 0; i < engine.WorkerCount; i++ {
		engine.createWorker(out, engine.Scheduler)
	}

	for _, request := range seeds {
		engine.Scheduler.Submit(request)
	}

	result := <-out
	for _, item := range result.Items {
		go func(item interface{}) {
			if liveData, ok := item.(model.LiveData); ok {
				name := config.RoomList[liveData.RoomId]
				switch liveData.LiveStatus {
				case 0:
					setRoomStatusFalse(liveData.RoomId)
				case 1:
					if !config.RoomStatusList[liveData.RoomId] {
						api.SendBarkMessage(name, "开播啦!")
						api.SendQQGroupMessage(config.GroupId, name+"开播啦!")
					}
					setRoomStatusTrue(liveData.RoomId)
				case 2:
					setRoomStatusFalse(liveData.RoomId)
				}
			}
			// TODO: 将数据放入数据库
			engine.ItemChan <- item
		}(item)
	}
}

// createWorker Worker创建/**
func (engine *ConcurrentEngine) createWorker(
	out chan Result, ready Scheduler) {
	go func() {
		// FIXME: 协程未关闭
		// for {
			in := ready.WorkerChan()
			ready.WorkerReady(in)
			request := <-in
			result, err := engine.RequestProcessor(request)
			if err == nil {
				// continue
				out <- result
			}
		// }
	}()
}

func setRoomStatusFalse(roomId int) {
	if config.RoomStatusList[roomId] {
		config.RoomStatusList[roomId] = false
	}
}

func setRoomStatusTrue(roomId int) {
	if !config.RoomStatusList[roomId] {
		config.RoomStatusList[roomId] = true
	}
}
