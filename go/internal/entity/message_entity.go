package entity

import (
	"encoding/json"
)


type Message struct {
	OS 			string			`json:"os"`
	Architecture string			`json:"arch"`
	CPUs		int 			`json:"cpus"`
	Threads		int 			`json:"threads"`
	GoRoutine 	int				`json:"go_routine"`
	MemoryAllocated uint64 		`json:"memory_allocated"`
	TotalMemoryAllocated uint64 	`json:"total_memory_allocated"`
	MemoryObtainedSystem uint64 	`json:"memory_obtained_system"`
	GarbaceCollection uint32 		`json:"garbage_collection"`
	GoVersion		string		`json:"go_version"`
}

func (message *Message) ToJSON() ([]byte,error) {
	return json.MarshalIndent(message,"","	")
}

func JSONToMessage(data []byte) (*Message,error) {
	var message Message
	err := json.Unmarshal(data,&message)
	return &message,err
}