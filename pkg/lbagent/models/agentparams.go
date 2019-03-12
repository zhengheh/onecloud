package models

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"text/template"

	"yunion.io/x/pkg/utils"

	"yunion.io/x/onecloud/pkg/mcclient/models"
)

func dataFromParams(p interface{}) map[string]interface{} {
	rv := reflect.ValueOf(p)
	if rv.Kind() != reflect.Struct {
		panic(fmt.Sprintf("unexpected kind: %#v", p))
	}
	rt := rv.Type()

	r := map[string]interface{}{}
	for i := 0; i < rv.NumField(); i++ {
		f := rt.Field(i)
		fn := utils.CamelSplit(f.Name, "_")
		if fn == "" {
			continue
		}
		v := rv.Field(i)
		if !v.IsValid() {
			continue
		}
		r[fn] = v.Interface()
	}
	return r
}

type AgentParams struct {
	AgentModel           *models.LoadbalancerAgent
	KeepalivedConfigTmpl *template.Template
	HaproxyConfigTmpl    *template.Template
	TelegrafConfigTmpl   *template.Template
	Data                 map[string]map[string]interface{}
}

func NewAgentParams(agent *models.LoadbalancerAgent) (*AgentParams, error) {
	b64s := map[string]string{
		"keepalived_conf_tmpl": agent.Params.KeepalivedConfTmpl,
		"haproxy_conf_tmpl":    agent.Params.HaproxyConfTmpl,
		"telegraf_conf_tmpl":   agent.Params.TelegrafConfTmpl,
	}
	tmpls := map[string]*template.Template{}
	for name, b64 := range b64s {
		d, err := base64.StdEncoding.DecodeString(b64)
		if err != nil {
			return nil, fmt.Errorf("%s: invalid base64 string: %s", name, err)
		}
		tmpl, err := template.New(name).Parse(string(d))
		if err != nil {
			return nil, fmt.Errorf("%s: invalid template: %s", name, err)
		}
		tmpls[name] = tmpl
	}
	dataAgent := map[string]interface{}{
		"id":   agent.Id,
		"name": agent.Name,
	}
	data := map[string]map[string]interface{}{
		"agent":    dataAgent,
		"vrrp":     dataFromParams(agent.Params.Vrrp),
		"haproxy":  dataFromParams(agent.Params.Haproxy),
		"telegraf": dataFromParams(agent.Params.Telegraf),
	}
	agentParams := &AgentParams{
		AgentModel:           agent,
		KeepalivedConfigTmpl: tmpls["keepalived_conf_tmpl"],
		HaproxyConfigTmpl:    tmpls["haproxy_conf_tmpl"],
		TelegrafConfigTmpl:   tmpls["telegraf_conf_tmpl"],
		Data:                 data,
	}
	return agentParams, nil
}

func (p *AgentParams) Equal(p2 *AgentParams) bool {
	if p == nil && p2 == nil {
		return true
	}
	if p == nil || p2 == nil {
		return false
	}
	agentP := p.AgentModel
	agentP2 := p2.AgentModel
	if agentP.Params != agentP2.Params {
		return false
	}
	return true
}

func (p *AgentParams) setXxParams(xx, k string, v interface{}) map[string]interface{} {
	var dt map[string]interface{}
	d, ok := p.Data[xx]
	if !ok {
		dt = map[string]interface{}{}
		p.Data[xx] = dt
	} else {
		dt = d
	}
	dt[k] = v
	return dt
}

func (p *AgentParams) SetVrrpParams(k string, v interface{}) map[string]interface{} {
	return p.setXxParams("vrrp", k, v)
}

func (p *AgentParams) SetHaproxyParams(k string, v interface{}) map[string]interface{} {
	return p.setXxParams("haproxy", k, v)
}

func (p *AgentParams) SetTelegrafParams(k string, v interface{}) map[string]interface{} {
	return p.setXxParams("telegraf", k, v)
}

func (p *AgentParams) KeepalivedConfig() {
}
