package cherryCode

const (
	OK                      = 0
	SessionUIDNotBind       = 10
	DiscoveryNotFoundNode   = 11
	NodeRequestError        = 12
	RPCNetError             = 20
	PRCUnmarshalError       = 21
	RPCMarshalError         = 22
	RPCRemoteExecuteError   = 23
	ActorPathIsNil          = 24
	ActorFuncNameError      = 25
	ActorConvertPathError   = 26
	ActorMarshalError       = 27
	ActorUnmarshalError     = 28
	ActorCallFail           = 29
	ActorSourceEqualTarget  = 30
	ActorPublishRemoteError = 31
	ActorChildIDNotFound    = 32
	ActorCallTimeout        = 33
	ActorIDIsNil            = 34
)

func IsOK(code int32) bool {
	return code == OK
}

func IsFail(code int32) bool {
	return code != OK
}
