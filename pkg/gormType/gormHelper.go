package gormType

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"math"
	"reflect"
	"strings"
	"time"
)

type Page struct {
	PageIndex int         `json:"page_index"`
	PageSize  int         `json:"page_size"`
	Total     int64       `json:"total"`
	PageNums  int64       `json:"page_nums"`
	List      interface{} `json:"list,omitempty"`
}

func NewMySqlDb(cfg string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg), &gorm.Config{})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to connect the ticket database,err: %s", err.Error()))
	}
	mdb, err2 := db.DB()
	if err2 != nil {
		return nil, errors.New(fmt.Sprintf("failed to new a ticket database, err: %s", err.Error()))
	}
	err = mdb.Ping()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to ping the ticket database,err: %s", err.Error()))
	}
	mdb.SetMaxOpenConns(500)
	mdb.SetConnMaxLifetime(time.Hour)
	return db, nil
}

func Close(db *gorm.DB) error {
	mdb, err := db.DB()
	if err != nil {
		return err
	}
	return mdb.Close()
}

func Count(db *gorm.DB, where interface{}, model interface{}) (int64, error) {
	tx := db.Model(model)
	tx, e := Where(tx, where)
	if e != nil {
		return 0, e
	}
	var count int64
	e = tx.Count(&count).Error
	if e != nil {
		return 0, e
	}
	return count, nil
}

func Get[T any](db *gorm.DB, where interface{}) (*T, error) {
	var entity, output T
	tx := db.Model(entity)
	tx, e := Where(tx, where)
	if e != nil {
		return nil, e
	}
	e = tx.Take(&output).Error
	if e == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if e != nil {
		return nil, e
	}
	return &output, nil
}

func Where(db *gorm.DB, where interface{}) (*gorm.DB, error) {
	var err error
	t := reflect.TypeOf(where).Kind()
	if t == reflect.Struct || t == reflect.Map {
		db = db.Where(where)
	} else if t == reflect.Slice {
		for _, item := range where.([]interface{}) {
			item := item.([]interface{})
			column := item[0]
			if reflect.TypeOf(column).Kind() == reflect.String {
				count := len(item)
				if count == 1 {
					return nil, errors.New("length less than 2")
				}
				columnStr := column.(string)
				// 拼接参数形式
				if strings.Index(columnStr, "?") > -1 {
					db = db.Where(column, item[1:]...)
				} else {
					cond := "and" //cond
					opt := "="
					_opt := " = "
					var val interface{}
					if count == 2 {
						opt = "="
						val = item[1]
					} else {
						opt = strings.ToLower(item[1].(string))
						_opt = " " + strings.ReplaceAll(opt, " ", "") + " "
						val = item[2]
					}

					if count == 4 {
						cond = strings.ToLower(strings.ReplaceAll(item[3].(string), " ", ""))
					}

					/*
					   '=', '<', '>', '<=', '>=', '<>', '!=', '<=>',
					   'like', 'like binary', 'not like', 'ilike',
					   '&', '|', '^', '<<', '>>',
					   'rlike', 'regexp', 'not regexp',
					   '~', '~*', '!~', '!~*', 'similar to',
					   'not similar to', 'not ilike', '~~*', '!~~*',
					*/

					if strings.Index(" in notin ", _opt) > -1 {
						// val 是数组类型
						column = columnStr + " " + opt + " (?)"
					} else if strings.Index(" = < > <= >= <> != <=> like likebinary notlike ilike rlike regexp notregexp", _opt) > -1 {
						column = columnStr + " " + opt + " ?"
					}

					if cond == "and" {
						db = db.Where(column, val)
					} else {
						db = db.Or(column, val)
					}
				}
			} else if t == reflect.Map /*Map*/ {
				db = db.Where(item)
			} else {
				/*
					// 解决and 与 or 混合查询，但这种写法有问题，会抛出 invalid query condition
					db = db.Where(func(db *gorm.DB) *gorm.DB {
						db, err = BuildWhere(db, item)
						if err != nil {
							panic(err)
						}
						return db
					})*/

				db, err = Where(db, item)
				if err != nil {
					return nil, err
				}
			}
		}
	} else {
		return nil, errors.New("param error")
	}
	return db, nil
}

func Paging(db *gorm.DB, page *Page, condition interface{}, model interface{}) error {
	tx := db.Model(model)
	tx, e := Where(tx, condition)
	if e != nil {
		return e
	}
	tx.Count(&page.Total)
	if page.Total == 0 {
		page.List = []interface{}{}
		return nil
	}
	//t := reflect.TypeOf(model)
	//list := reflect.Zero(reflect.SliceOf(t)).Interface().([]User) //failed to compile
	//list := reflect.New(reflect.ArrayOf(0, t)).Interface().(*[0]User)

	e = tx.Scopes(Paginate(page)).Find(&page.List).Error
	if e != nil {
		return e
	}
	return nil
}

func Tx(db *gorm.DB, funcArr ...func(db *gorm.DB) error) (err error) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = fmt.Errorf("%v", r)
		}
	}()
	for _, f := range funcArr {
		err = f(tx)
		if err != nil {
			tx.Rollback()
			return
		}
	}
	err = tx.Commit().Error
	return err
}

func Paginate(page *Page) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page.PageIndex < 1 {
			page.PageIndex = 1
		}
		if page.PageSize < 1 {
			page.PageSize = 10
		}
		page.PageNums = int64(math.Ceil(float64(page.Total) / float64(page.PageSize)))
		return db.Offset((page.PageIndex - 1) * page.PageSize).Limit(page.PageSize)
	}
}

func Status(status string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status in (?)", strings.Split(status, ","))
	}
}

func ConvertToMap(updateFields []string, obj interface{}) map[string]interface{} {
	data := make(map[string]interface{})
	values := reflect.ValueOf(obj)
	t := reflect.TypeOf(obj)
	for i := 0; i < t.NumField(); i++ {
		fName := t.Field(i).Name
		for _, updateField := range updateFields {
			if fName == updateField {

				asciiCodes := []byte(updateField)
				var lowAsciiCodes []byte
				for j, aCode := range asciiCodes {
					if aCode >= 65 && aCode <= 90 {
						if j > 0 {
							lowAsciiCodes = append(lowAsciiCodes, 95) //下划线
						}
						lowAsciiCodes = append(lowAsciiCodes, aCode+32)
						continue
					}
					lowAsciiCodes = append(lowAsciiCodes, aCode)
				}
				data[string(lowAsciiCodes)] = values.FieldByName(fName).Interface()
			}
		}
	}
	return data
}
