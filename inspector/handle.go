package main

import (
	"owl/common/tcp"
	"owl/common/types"
)

type InspectorHandle struct {
}

func (this *InspectorHandle) MakeSession(sess *tcp.Session) {
	lg.Info("%s new connection ", sess.RemoteAddr())
}

func (this *InspectorHandle) LostSession(sess *tcp.Session) {
	lg.Info("%s disconnect ", sess.RemoteAddr())
}

func (this *InspectorHandle) Handle(sess *tcp.Session, data []byte) {
	defer func() {
		if err := recover(); err != nil {
			lg.Error("Recovered in HandleMessage", err)
		}
	}()
	mt := types.AlarmMessageType(data[0])
	lg.Debug("Receive %v %v", types.AlarmMessageTypeText[mt], string(data[1:]))
	switch mt {
	case types.ALAR_MESS_GET_INSPECTOR_TASK_RESP:
		get_task_resp := &types.GetTasksResp{}
		if err := get_task_resp.Decode(data[1:]); err != nil {
			lg.Error(err.Error())
			return
		}
		inspector.taskPool.PutTasks(get_task_resp.AlarmTasks)
	default:
		lg.Error("Unknown option: %v", mt)
	}
}
