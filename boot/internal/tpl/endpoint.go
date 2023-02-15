package tpl

var Endpoint = `package ${useCasePkg}

import (
	"context"
	"github.com/zander-84/seagull/contrib/storage"
	"github.com/zander-84/seagull/endpoint"
	"github.com/zander-84/seagull/endpoint/wraptransporter"
	"${project}/apps/${server}/internal/entity"
	"${project}/apps/${server}/internal/usecase"
)

func Get${EntityName}(ctx context.Context, request any) (response any, err error) {
	in := request.(*Get${EntityName}Codec)
	return usecase.${EntityName}UsesCase.Get${EntityName}(context.Background(), in.Id)
}

func Create${EntityName}(ctx context.Context, request any) (response any, err error) {
	in := request.(*Create${EntityName}Codec)
	${entityName} := new(entity.${EntityName})
	// todo check assign
	${assignCreateFields}
	err = usecase.${EntityName}UsesCase.Create${EntityName}(context.Background(), ${entityName})
	return ${entityName}, err
}

func Update${EntityName}(ctx context.Context, request any) (response any, err error) {
	in := request.(*Update${EntityName}Codec)

	${entityName}, err := usecase.${EntityName}UsesCase.Get${EntityName}(context.Background(), in.Id)
	if err != nil {
		return nil, err
	}
		// todo check assign
	${assignUpdateFields}

	err = usecase.${EntityName}UsesCase.Update${EntityName}(context.Background(), in.Id, in.Version, ${entityName})
	return ${entityName}, err
}

func BatchGet${EntityName}(ctx context.Context, request any) (response any, err error) {
	in := request.(*Batch${EntityName}Codec)

	data, err := usecase.${EntityName}UsesCase.BatchGet${EntityName}(context.Background(), in.Ids)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func Search${EntityName}(ctx context.Context, request any) (response any, err error) {
	// in := request.(*Search${EntityName}Codec)
	transporter := endpoint.GetTransporter(ctx)
	searchMeta := storage.NewSearchMeta()
	searchMeta.SetPage(wraptransporter.GetPage(transporter)).SetPageSize(wraptransporter.GetPageSize(transporter))

	sqlBuilder := storage.NewMysqlBuilder()
	//todo add where
	

	data, cnt, err := usecase.${EntityName}UsesCase.Search${EntityName}(context.Background(), searchMeta, sqlBuilder)
	if err != nil {
		return nil, err
	}
	wraptransporter.SetCount(transporter, cnt)

	return data, nil
}
`

var EndpointCodec = `package ${useCasePkg}

import (
	"context"
	"github.com/zander-84/seagull/endpoint"
	"github.com/zander-84/seagull/endpoint/wraptransporter"
	"github.com/zander-84/seagull/think"
	"github.com/zander-84/seagull/tool/conv"
	"github.com/zander-84/seagull/transport/http"
	"${project}/apps/${server}/internal/entity"
)

type Get${EntityName}Codec struct {
	${IdWithType}
}

func (Get${EntityName}Codec) HttpGetDecode(ctx context.Context, request any) (any, error) {
	out := new(Get${EntityName}Codec)
	httpCtx := ctx.(http.Context)
	${outId}
	return out, nil
}

func (Get${EntityName}Codec) HttpGetEncode(ctx context.Context, request any) (any, error) {
	httpCtx := ctx.(http.Context)
	// todo reset response 

	resp := think.NewSuccessResp(request)
	return httpCtx.JSON(resp.Code.HttpCode(), resp), nil
}

type Create${EntityName}Codec struct {
	${fields}
}

func (Create${EntityName}Codec) HttpPostDecode(ctx context.Context, request any) (any, error) {
	out := new(Create${EntityName}Codec)
	httpCtx := ctx.(http.Context)
	if err := httpCtx.BindJson(out); err != nil {
		return nil, think.ErrAlert(err.Error())
	}
	// todo validate
	return out, nil
}

func (Create${EntityName}Codec) HttpPostEncode(ctx context.Context, request any) (any, error) {
	httpCtx := ctx.(http.Context)

	// todo reset response
	resp := think.NewSuccessResp(request)

	return httpCtx.JSON(resp.Code.HttpCode(), resp), nil
}

type Update${EntityName}Codec struct {
	${IdWithType}
	${fields}
}

func (Update${EntityName}Codec) HttpPutDecode(ctx context.Context, request any) (any, error) {
	out := new(Update${EntityName}Codec)
	httpCtx := ctx.(http.Context)
	if err := httpCtx.BindJson(out); err != nil {
		return nil, think.ErrAlert(err.Error())
	}
	${outId}
	//todo validate
	return out, nil
}

func (Update${EntityName}Codec) HttpPutEncode(ctx context.Context, request any) (any, error) {
	httpCtx := ctx.(http.Context)
	// todo reset response
	resp := think.NewSuccessResp(request)
	return httpCtx.JSON(resp.Code.HttpCode(), resp), nil
}


type Batch${EntityName}Codec struct {
	${idsWithType}
}

func (Batch${EntityName}Codec) HttpGetDecode(ctx context.Context, request any) (any, error) {
	out := new(Batch${EntityName}Codec)
	httpCtx := ctx.(http.Context)
	if err := httpCtx.Bind(out); err != nil {
		return nil, think.ErrAlert(err.Error())
	}

	return out, nil
}

func (Batch${EntityName}Codec) HttpGetEncode(ctx context.Context, request any) (any, error) {

	httpCtx := ctx.(http.Context)
	${entityName}s, ok := request.([]entity.${EntityName})
	if !ok {
		return nil, think.ErrSystemSpace("request err")
	}
	resp := think.NewSuccessResp(request)
	return httpCtx.JSON(resp.Code.HttpCode(), ${entityName}s), nil
}


type Search${EntityName}Codec struct {
}

func (Search${EntityName}Codec) HttpGetDecode(ctx context.Context, request any) (any, error) {
	out := new(Search${EntityName}Codec)
	httpCtx := ctx.(http.Context)
	if err := httpCtx.Bind(out); err != nil {
		return nil, think.ErrAlert(err.Error())
	}

	return out, nil
}

func (Search${EntityName}Codec) HttpGetEncode(ctx context.Context, request any) (any, error) {
	transporter := endpoint.GetTransporter(ctx)

	httpCtx := ctx.(http.Context)
	${entityName}s, ok := request.([]entity.${EntityName})
	if !ok {
		return nil, think.ErrSystemSpace("request err")
	}
	resp := think.NewSuccessResp(request)
	return httpCtx.JSON(resp.Code.HttpCode(), map[string]any{
		"cnt":  wraptransporter.GetCount(transporter),
		"data": ${entityName}s,
	}), nil
}

`
