package template

import (
	_ "embed"
	"fmt"
	"github.com/fatih/color"
	"github.com/iyear/E5SubBot/app/bot/internal/model"
	"github.com/iyear/E5SubBot/pkg/utils"
	iso6391 "github.com/iyear/iso-639-1"
	"github.com/spf13/viper"
	tele "gopkg.in/telebot.v3"
	"io/fs"
	"path/filepath"
	"reflect"
	"strconv"
	"text/template"
	"time"
)

func Init(tmplPath string) error {
	color.Blue("loading templates...\n")
	paths, err := walkTemplates(tmplPath)
	if err != nil {
		return err
	}

	for _, p := range paths {
		t := model.Template{}
		if err = bindTmpl(p, &t); err != nil {
			return err
		}

		code := utils.GetFileName(p)
		setButtons(code, &t)
		color.Blue("\t- %s\n", code)
		ints[code] = &t
		codes = append(codes, code)
	}

	return nil
}

func tmplHook() viper.DecoderConfigOption {
	return viper.DecodeHook(func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		// bind template
		if f.Kind() == reflect.String && t == reflect.TypeOf(model.Tmpl{}) {
			tmpl, err := template.New(strconv.FormatInt(time.Now().UnixNano(), 10)).Parse(data.(string))
			if err != nil {
				return data, nil
			}
			return model.NewTmpl(tmpl), nil
		}

		// 不处理，循环内手动处理
		if t == reflect.TypeOf([][]tele.InlineButton{}) {
			return nil, nil
		}
		return data, nil
	})
}

func From(lang string) *model.TMPL {
	return &model.TMPL{
		I: iFrom(lang),
		B: bFrom(lang),
	}
}

func walkTemplates(dirPath string) ([]string, error) {
	paths := make([]string, 0)
	err := filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("unable to walk template path: %s, err:%v", path, err)
		}
		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		name := utils.GetFileName(path)
		if ext != ".toml" || !iso6391.ValidCode(name) {
			return fmt.Errorf("invalid template file: %s.Please check extension or name of the file", path)
		}

		paths = append(paths, path)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return paths, nil
}

func bindTmpl(path string, value interface{}) error {
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return err
	}
	if err := v.Unmarshal(value, tmplHook()); err != nil {
		return err
	}
	return nil
}
