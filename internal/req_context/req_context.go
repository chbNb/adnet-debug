package req_context

type ReqContext struct {
	serverIpUrl string
	serverIP    string
	remoteIP    string
	BucketSpeed int
	PreTokens   int
}

var ctx *ReqContext

func GetInstance() *ReqContext {
	if ctx == nil {
		ctx = &ReqContext{}
	}
	return ctx
}

func (ctx *ReqContext) UpdateServerIpInfo() {
	if len(ctx.serverIP) == 0 && len(ctx.serverIpUrl) > 0 {
		serverIp, err := GetServerIP(ctx.serverIpUrl)
		if err == nil {
			ctx.SetServerIp(serverIp)
		}
	}
	if len(ctx.remoteIP) == 0 {
		remoteIp, err := GetInternalIP()
		if err == nil {
			ctx.SetRemoteIp(remoteIp)
		}
	}
}

func (ctx *ReqContext) SetServerIpUrl(url string) {
	ctx.serverIpUrl = url
}

func (ctx *ReqContext) SetServerIp(serverIp string) {
	ctx.serverIP = serverIp
}

func (ctx *ReqContext) SetRemoteIp(remoteIp string) {
	ctx.remoteIP = remoteIp
}

func (ctx *ReqContext) ServerIP() string {
	return ctx.serverIP
}

func (ctx *ReqContext) RemoteIP() string {
	return ctx.remoteIP
}
