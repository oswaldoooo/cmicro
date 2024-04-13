package delaymap

type DelayMap[Key any] interface {
	SetCallBackWhenExpire(Key, func())
}
