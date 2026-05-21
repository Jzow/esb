package util

import (
	"context"
	"errors"
	"fmt"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"github.com/olivere/elastic/v7"
	"math"
	"reflect"
	"strings"
	"time"
)

var ReportEsClient *elastic.Client

func SetupElastic() error {
	beginTime := time.Now()
	var err error
	//elastic.SetBasicAuth("user", "secret")
	ReportEsClient, err = elastic.NewClient(elastic.SetURL(setting.CommonEsSetting.Host),
		elastic.SetBasicAuth(setting.CommonEsSetting.Account, setting.CommonEsSetting.Password),
		elastic.SetSniff(false),
		elastic.SetTraceLog(EsLogger{}))
	Log("[SetupElastic]connect to elasticsearch,cost %v", time.Since(beginTime), err)
	if err != nil {
		return err
	}
	return nil
}

type EsLogger struct{}

func (EsLogger) Printf(format string, v ...interface{}) {
	logging.Infof(format, v)
}

// EsSearchReq T:ESDoc解析的强类型 Fields:ES查询的字段 FilterFields:ES返回的筛选字段
type EsSearchReq[T any] struct {
	*Paging[T]
	Index        string
	Sorts        map[string]bool
	Fields       map[string]any
	FilterFields []string
}

func ElasticSearch[T any](form *EsSearchReq[T]) error {
	logTag := fmt.Sprintf("[ElasticSearch][%s]", time.Now().Format("20060102150405"))
	if form == nil || len(form.Fields) == 0 {
		logging.Infof("%s invalid param", logTag)
		return errors.New("invalid param")
	}

	if form.Paging == nil {
		form.Paging = &Paging[T]{PageIndex: 1, PageSize: 10000}
	}

	form.List = []T{}
	if form.PageIndex < 1 {
		form.PageIndex = 1
	}
	if form.PageSize < 1 {
		form.PageSize = 10
	}
	logging.Infof("%s PageIndex=%d,PageSize=%d", logTag, form.PageIndex, form.PageSize)

	if ReportEsClient == nil {
		err := SetupElastic()
		if err != nil {
			logging.Infof("%s failed to connect the elastic client,error:%s", logTag, err.Error())
			return nil
		}
	}

	query := elastic.NewBoolQuery()
	for field, valObj := range form.Fields {
		switch reflect.TypeOf(valObj).Kind() {
		case reflect.Slice, reflect.Array:
			array, _ := To[[]any](valObj)
			query.Filter(elastic.NewTermsQueryFromStrings(field, ToStringArray(*array, ",")...))
			break
		case reflect.Map:
			if valMap, ok := valObj.(map[string]any); ok {
				if field == "$or" {
					var or []elastic.Query
					for k, v := range valMap {
						or = append(or, elastic.NewTermQuery(k, v))
					}
					query.Should(or...).MinimumNumberShouldMatch(1) //至少一个条件成立
				} else {
					if valType, ok := valMap["type"]; ok && valType == "SecondTimeStamp" {
						if val, ok := valMap["value"]; ok {
							query.Filter(GetSecondTimeStampRangeQuery(field, ObjToStr(val)))
						}
					}
				}
			}
			break
		default:
			query.Filter(elastic.NewTermQuery(field, valObj))
		}
	}

	var sorts []elastic.Sorter
	for field, isAscSort := range form.Sorts {
		sorts = append(sorts, elastic.NewFieldSort(field).Order(isAscSort))
	}
	srv := ReportEsClient.Search().Index(form.Index).Query(query).SortBy(sorts...).From((form.PageIndex - 1) * form.PageSize).Size(form.PageSize)
	if len(form.FilterFields) > 0 {
		srv = srv.FetchSourceContext(elastic.NewFetchSourceContext(true).Include(form.FilterFields...))
	}
	searchResult, err := srv.Do(context.TODO())
	if err != nil {
		logging.Errorf("%s failed to search, error:%s", logTag, err.Error())
		return errors.New("failed to search")
	}
	if searchResult == nil || searchResult.Hits == nil {
		logging.Infof("%s search empty result", logTag)
		return nil
	}
	form.Total = searchResult.TotalHits()
	form.PageNums = int64(math.Ceil(float64(form.Total) / float64(form.PageSize)))
	logging.Infof("%s search es result count: %d/%d", logTag, len(searchResult.Hits.Hits), form.Total)

	if searchResult.Hits == nil || searchResult.Hits.Hits == nil {
		return nil
	}
	for _, hit := range searchResult.Hits.Hits {
		if hit.Source != nil {
			t, _ := UnmarshalObj[T](hit.Source)
			form.List = append(form.List, t)
		}
	}

	/*	var t T
		for _, item := range searchResult.Each(reflect.TypeOf(t)) {
			form.List = append(form.List, item.(T))
		}*/

	return nil
}

func GetSecondTimeStampRangeQuery(fieldName, fieldValue string) *elastic.RangeQuery {
	rangeQuery := elastic.NewRangeQuery(fieldName).Format("epoch_second").Lte("now")
	if fieldValue != "" {
		times := strings.Split(fieldValue, ",")
		startTime := ObjToUint(times[0])
		if startTime == 0 {
			startTime = GetTimeStampByStrFormat(times[0])
		}
		if startTime > 0 {
			rangeQuery.Gte(startTime)
			if len(times) > 1 {
				endTime := ObjToUint(times[1])
				if endTime == 0 {
					endTime = GetTimeStampByStrFormat(times[1])
				}
				if endTime > 0 {
					rangeQuery.Lte(endTime)
				}
			}
		}
	}
	return rangeQuery
}
