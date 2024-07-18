##api调用接口

对接第三方应用时，配置中使用
```
api:
  user:
    # 调用下游的服务名称
    service: user
    # 请求完整地址
    domain: http://10.20.22.10:9051
    # 超时配置，time.Duration 类型
    timeout: 5000ms
    # 重试次数，最多执行retry+1次
    retry: 2
```
代码调用式例如下
```
func SubmitTask(ctx *gin.Context, taskType string, SummaryId string, input string) (CommitResp, error) {
	infer := CommitResp{}
	apiReq := ServingPostCommitReq{
		SessionId:  SummaryId,
		TaskType:   taskType,
		Input:      input,
		CommitType: CommitType{IsWs: true},
	}
	resp, err := conf.Resource.Api["jobd"].HttpPost(ctx, "/jobd/committer/Commit", api.HttpRequestOptions{
		RequestBody: apiReq,
		Encode:      api.EncodeJson,
	})
	if err != nil {
		return infer, err
	}

	if resp.HttpCode != 200 {
		return infer, err
	}
	if err = json.Unmarshal(resp.Response, &infer); err != nil {
		return infer, err
	}
	return infer, nil
}
```