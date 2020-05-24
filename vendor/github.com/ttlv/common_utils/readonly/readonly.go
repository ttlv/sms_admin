package readonly

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/qor/admin"
	"github.com/qor/qor/utils"
)

var disabledExportResources []*admin.Resource

func Setup(adm *admin.Admin, cfg ...interface{}) {
	adm.RegisterViewPath("github.com/ttlv/common_utils/readonly/views")
	for _, res := range adm.GetResources() {
		if res.GetTheme("readonly") != nil {
			adm.GetRouter().Get(fmt.Sprintf("/%v/export", res.ToParam()), func(c *admin.Context) {
				res := c.Resource
				if c.Request.URL.Query()["timestamp"] == nil {
					c.Writer.WriteHeader(500)
					c.Writer.Write([]byte(`timestamp不能为空`))
					return
				}
				maxCount := exportMaxCount(cfg, c)
				c.PerPage(maxCount + 1)
				results, err := c.FindMany()
				if err != nil {
					c.Writer.Write([]byte(err.Error()))
					return
				}
				size := reflect.Indirect(reflect.ValueOf(results)).Len()
				if size > maxCount {
					c.Writer.WriteHeader(500)
					c.Writer.Write([]byte(fmt.Sprintf(`下载结果超过%v, 请过滤数据再下载.`, maxCount)))
					return
				}
				timestamp := c.Request.URL.Query()["timestamp"][0]
				cookie := http.Cookie{Name: "fileDownload", Value: timestamp}
				http.SetCookie(c.Writer, &cookie)

				export := exportHeaders(c, res) + "\n"
				export += exportRecords(c, res, results)
				c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename="+"%v_%v.csv", strings.Replace(utils.ToParamString(c.Resource.ToParam()), "readonly_", "", 1), time.Now().UnixNano()))
				c.Writer.Header().Set("Content-Type", c.Request.Header.Get("Content-Type"))
				c.Writer.Write([]byte(export))
			})
		}
	}
	adm.RegisterFuncMap("is_not_last", func(index int, count int) bool {
		return count-1 != index
	})
	adm.RegisterFuncMap("enable_export", func(resource *admin.Resource) bool {
		for _, res := range disabledExportResources {
			if res == resource {
				return false
			}
		}
		return true
	})
}

func DisableExport(res *admin.Resource) {
	disabledExportResources = append(disabledExportResources, res)
}

func exportHeaders(c *admin.Context, res *admin.Resource) string {
	headers := []string{}
	for _, section := range res.IndexAttrs() {
		for _, rows := range section.Rows {
			meta := res.GetMeta(rows[0])
			key := fmt.Sprintf("%v.attributes.%v", res.ToParam(), meta.Label)
			headers = append(headers, string(c.Admin.T(c.Context, key, meta.Label)))
		}
	}
	return strings.Join(headers, ",")
}

func exportRecords(c *admin.Context, res *admin.Resource, results interface{}) string {
	resultValues := reflect.Indirect(reflect.ValueOf(results))
	exportValues := []string{}
	for i := 0; i < resultValues.Len(); i++ {
		vv := []string{}
		for _, section := range res.IndexAttrs() {
			for _, rows := range section.Rows {
				var (
					maxGoroutines = 10
					guard         = make(chan struct{}, maxGoroutines)
					wg            sync.WaitGroup
				)
				wg.Add(1)
				guard <- struct{}{}
				go func(meta *admin.Meta) {
					vv = append(vv, fmt.Sprintf("%v", meta.GetFormattedValuer()(resultValues.Index(i).Interface(), c.Context)))
					wg.Done()
					<-guard
				}(res.GetMeta(rows[0]))
				wg.Wait()
			}
		}
		exportValues = append(exportValues, strings.Join(vv, ","))
	}
	return strings.Join(exportValues, "\n")
}

func exportMaxCount(cfg []interface{}, c *admin.Context) int {
	if c.Request.URL.Query()["max"] != nil {
		if c, err := strconv.Atoi(c.Request.URL.Query()["max"][0]); err == nil {
			if c < 1000000 {
				return c
			}
		}
	}
	if len(cfg) != 0 {
		if m, ok := cfg[0].(int); ok {
			return m
		}
	}
	return 100000
}
