package client

import (
	"io/ioutil"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/vitelabs/vite-portal/shared/pkg/interfaces"
	"github.com/vitelabs/vite-portal/shared/pkg/types"
	"github.com/vitelabs/vite-portal/shared/pkg/util/httputil"
	"github.com/vitelabs/vite-portal/shared/pkg/util/jsonutil"
)

type ViteClient struct {
	url       string
	requestId uint64
	client    *retryablehttp.Client
}

func NewViteClient(url string) *ViteClient {
	c := &ViteClient{
		url:       url,
		requestId: 0,
		client:    retryablehttp.NewClient(),
	}
	c.client.RetryWaitMin = 500 * time.Millisecond
	c.client.RetryWaitMax = 3 * time.Second
	c.client.RetryMax = 3
	return c
}

func (c *ViteClient) createRequest(method string, params []interface{}) *types.RpcRequest {
	atomic.AddUint64(&c.requestId, 1)
	return &types.RpcRequest{
		Id:      int(c.requestId),
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
	}
}

func (c *ViteClient) Send(method string, params []interface{}, v interfaces.RpcResponseI) error {
	request := c.createRequest(method, params)
	requestBody, err := jsonutil.ToByte(request)
	if err != nil {
		return err
	}
	r, err := c.client.Post(c.url, httputil.ContentTypeJson, requestBody)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = httputil.ExtractModelFromBody(data, v)
	if err != nil {
		return err
	}
	err = v.GetError()
	if err != nil {
		return err
	}
	return nil
}

func (c *ViteClient) GetSnapshotChainHeight() (int64, error) {
	resp := types.RpcResponse[string]{}
	err := c.Send("ledger_getSnapshotChainHeight", nil, &resp)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(resp.Result, 10, 64)
}

func (c *ViteClient) GetLatestAccountBlock(addr string) (types.RpcViteLatestAccountBlockResponse, error) {
	resp := types.RpcResponse[types.RpcViteLatestAccountBlockResponse]{}
	err := c.Send("ledger_getLatestAccountBlock", []interface{}{addr}, &resp)
	if err != nil {
		return *new(types.RpcViteLatestAccountBlockResponse), err
	}
	return resp.Result, nil
}
