package engine

import (
	"api"
	"config"
	"fmt"
	"monitor/model"
)

// Run 多线程调度引擎/**
func (engine *ConcurrentEngine) Run(seeds ...Request) {

	out := make(chan Result)
	engine.Scheduler.Run()

	for i := 0; i < engine.WorkerCount; i++ {
		engine.createWorker(engine.Scheduler.WorkerChan(), out, engine.Scheduler)
	}

	for _, request := range seeds {
		engine.Scheduler.Submit(request)
	}

	result := <-out
	for _, item := range result.Items {
		go func(item interface{}) {
			if liveData, ok := item.(model.LiveData); ok {
				name := config.RoomList[liveData.RoomId]
				fmt.Print(name + ": ")
				switch liveData.LiveStatus {
				case 0:
					fmt.Println("尚未直播")
					setRoomStatusFalse(liveData.RoomId)
				case 1:
					fmt.Println("直播中")
					if !config.RoomStatusList[liveData.RoomId] {
						api.SendBarkMessage(name, "开播啦!")
						api.SendQQGroupMessage(config.GroupId, name+"开播啦!")
					}
					setRoomStatusTrue(liveData.RoomId)
				case 2:
					fmt.Println("轮播中")
					setRoomStatusFalse(liveData.RoomId)
				}
			}
			// TODO: 将数据放入数据库
			engine.ItemChan <- item
		}(item)
	}
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
