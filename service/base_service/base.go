package base_service

import (
	"encoding/json"
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/gormType"
	"gorm.io/gorm"
)

type DBService[T models.DBEntity] struct {
	Entity *T
	*gormType.Page
	*gorm.DB
}

func NewService[T models.DBEntity](entity *T, args ...any) *DBService[T] {
	var srv DBService[T]
	srv.DB = models.UDB
	srv.Entity = entity
	if len(args) > 0 {
		jBytes, err := json.Marshal(args[0])
		if err == nil {
			var page gormType.Page
			err = json.Unmarshal(jBytes, &page)
			if err == nil {
				page.List = []T{}
				srv.Page = &page
			}
		}
	}
	return &srv
}

func (t *DBService[T]) CountByField(where any) (int64, error) {
	var entity T
	return gormType.Count(t.DB, t.getFilterFields(where), entity)
}

func (t *DBService[T]) Get(where any) (*T, error) {
	return gormType.Get[T](t.DB, t.getFilterFields(where))
}

func (t *DBService[T]) UpdateFieldById(field []string, idField ...string) int64 {
	return t.Model(t.Entity).Updates(t.GetMaps(field)).RowsAffected
}

func (t *DBService[T]) Paging(where []string, other ...map[string]any) error {
	var entity T
	whereMap := t.GetMaps(where)
	for _, m := range other {
		for k, v := range m {
			whereMap[k] = v
		}
	}
	return gormType.Paging(t.DB, t.Page, whereMap, entity)
}

func (t *DBService[T]) GetAll(where []string, orders []string, other ...map[string]any) (list []T, err error) {
	var entity T
	whereMap := t.GetMaps(where)
	for _, m := range other {
		for k, v := range m {
			whereMap[k] = v
		}
	}
	tx := t.Model(entity).Where(whereMap)
	for _, order := range orders {
		tx.Order(order)
	}
	err = tx.Find(&list).Error
	return list, err
}

func (t *DBService[T]) getFilterFields(where any) any {
	switch where.(type) {
	case []string:
		return t.GetMaps(where.([]string))
	}
	return where
}

func (t *DBService[T]) GetMaps(where []string) map[string]interface{} {
	return gormType.ConvertToMap(where, *t.Entity)
}
