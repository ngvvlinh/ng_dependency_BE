package driver

import (
	"context"

	"o.o/api/etelecom/call_state"
	etelecomxserviceclient "o.o/backend/pkg/integration/telecom/etelecomxservice/client"
)

type Driver struct {
	client *etelecomxserviceclient.Client
}

func New(token string) *Driver {
	client := etelecomxserviceclient.New(token)
	return &Driver{
		client: client,
	}
}

func (d *Driver) GetCallLogs(ctx context.Context, req *GetCallLogsRequest) (res *GetCallLogsResponse, _ error) {
	getCallLogReq := &etelecomxserviceclient.GetCallLogsRequest{
		ScrollID: req.ScrollID,
	}
	if !req.StartedAt.IsZero() {
		getCallLogReq.StartTime = req.StartedAt.Unix()
	}
	if !req.EndedAt.IsZero() {
		getCallLogReq.EndTime = req.EndedAt.Unix()
	}

	getCallLogsResp, err := d.client.GetCallLogs(ctx, getCallLogReq)
	if err != nil {
		return nil, err
	}

	res = &GetCallLogsResponse{
		ScrollID: getCallLogsResp.ScrollID.String(),
	}

	for _, callLog := range getCallLogsResp.Sessions {
		tartgets := getCallTargets(callLog)
		callState := callLog.CallStatus.ToCallState()
		hotlineNumber := callLog.DidCid
		if hotlineNumber == "" {
			hotlineNumber = callLog.OutboundCallerID
		}
		callLogRes := &CallLog{
			CallID:        callLog.CallID.String(),
			CallStatus:    string(callLog.CallStatus),
			Caller:        callLog.Caller.String(),
			Callee:        callLog.Callee.String(),
			Direction:     callLog.Direction.String(),
			StartedAt:     callLog.StartTime.ToTime(),
			EndedAt:       callLog.EndedTime.ToTime(),
			Duration:      callLog.TalkDuration.Int(),
			CallTargets:   tartgets,
			CallState:     callState,
			HotlineNumber: hotlineNumber.String(),
			SessionID:     callLog.SessionID.String(),
		}
		fileRecording := callLog.RecordingFileURL.String()
		if fileRecording != "" {
			callLogRes.AudioURLs = append(callLogRes.AudioURLs, fileRecording)
		}
		res.CallLogs = append(res.CallLogs, callLogRes)
	}

	return res, nil
}

func getCallTargets(session *etelecomxserviceclient.SessionCallLog) (res []*CallTarget) {
	callTargetsMap := make(map[string]*CallTarget)

	for _, target := range session.CallTargets {
		targetNumber := target.TargetNumber.String()
		callState := target.Status.ToCallState()

		if _, ok := callTargetsMap[targetNumber]; !ok || callState == call_state.Answered {
			// make sure target number is unique
			callTargetsMap[targetNumber] = &CallTarget{
				TargetNumber: targetNumber,
				TalkDuration: target.TalkDuration.Int(),
				CallState:    callState,
				AnsweredTime: target.AnsweredTime.ToTime(),
				EndedTime:    target.EndedTime.ToTime(),
			}
		}
	}

	for _, target := range callTargetsMap {
		res = append(res, target)
	}
	return res
}
