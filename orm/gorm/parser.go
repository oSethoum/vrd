package gorm

import (
	"fmt"
	"strings"
	"vrd/types"
)

type Parser struct {
	e *Engine
}

func NewParser(e *Engine) *Parser {
	return &Parser{
		e: e,
	}
}

func (p *Parser) Start() {
	p.e.state.Models = make(map[string]Model)
	ignored := []string{"id", "created_at", "updated_at", "deleted_at"}
	if p.e.config.Gorm.GormModel {
		p.e.imports = append(p.e.imports, "time")
	}
	for _, t := range p.e.vuerd.TableState.Tables {
		model := Model{
			Name:      t.Name,
			Columns:   make(map[string]Column),
			GormModel: !p.e.h.Contains(p.e.h.Clean(t.Comment), "model:ignore"),
		}

		for _, c := range t.Columns {
			if (p.e.config.Gorm.GormModel || p.e.h.Contains(p.e.h.Clean(t.Comment), "model:ignore")) && p.e.h.InArray(ignored, c.Name) {
				continue
			}

			json := ""
			column := Column{
				Name: p.e.h.Pascal(c.Name),
				Type: GormTypesMap[types.VuerdTypes[p.e.h.Lower(c.DataType)]],
			}

			switch column.Type {
			case "time.Time":
				if !p.e.h.InArray(p.e.imports, "time") {
					p.e.imports = append(p.e.imports, "time")
				}
			case "datatypes.JSON":
				if !p.e.h.InArray(p.e.imports, "gorm.io/datatypes") {
					p.e.imports = append(p.e.imports, "gorm.io/datatypes")
				}
				model.JsonFields = append(model.JsonFields, p.e.h.SpecialCamel(c.Name))
			}

			if p.e.h.Contains(p.e.h.Clean(c.Comment), "json:ignore") {
				json = "json:\"-\""
			} else {
				json = fmt.Sprintf("json:\"%s,omitempty\"", p.e.h.SpecialCamel(c.Name))
			}

			gormOps := []string{}

			if c.Option.PrimaryKey {
				gormOps = append(gormOps, "primaryKey")
			}

			if c.Option.AutoIncrement {
				gormOps = append(gormOps, "autoIncrement")
			}

			if c.Option.NotNull {
				gormOps = append(gormOps, "not null")
			}

			if !p.e.h.Contains(c.Comment, "json:ignore") {
				column.Type = "*" + column.Type
			}

			if c.Option.Unique {
				gormOps = append(gormOps, "unique")
			}

			if len(c.Default) > 0 {
				gormOps = append(gormOps, fmt.Sprintf("default:%s", c.Default))
			}

			gorm := ""

			if p.e.h.Contains(p.e.h.Clean(c.Comment), "gorm:ignore") {
				gorm = "gorm:\"-\""
			} else if len(gormOps) > 0 {
				gorm = fmt.Sprintf("gorm:\"%s\"", strings.Join(gormOps, ", "))
			}
			column.Options = fmt.Sprintf("`%s %s`", json, gorm)
			model.Columns[c.Id] = column
		}

		p.e.state.Models[t.Id] = model
	}

	for _, r := range p.e.vuerd.RelationshipState.Relationships {
		start := p.e.state.Models[r.Start.TableId]
		end := p.e.state.Models[r.End.TableId]

		sName := p.e.h.Pascal(start.Name)
		end.Columns[sName] = Column{
			Name:    sName,
			Type:    fmt.Sprintf("*%s", start.Name),
			Options: fmt.Sprintf("`json:\"%s,omitempty\"`", p.e.h.SpecialCamel(start.Name)),
		}

		switch r.RelationshipType {
		case "ZeroN", "OneN":
			Name := p.e.h.Pascals(end.Name)
			start.Columns[Name] = Column{
				Name:    Name,
				Type:    fmt.Sprintf("[]%s", end.Name),
				Options: fmt.Sprintf("`json:\"%s,omitempty\"`", p.e.h.SpecialCamels(end.Name)),
			}
		case "OneOnly", "ZeroOne":
			eName := p.e.h.Pascal(end.Name)
			start.Columns[eName] = Column{
				Name:    eName,
				Type:    fmt.Sprintf("*%s", end.Name),
				Options: fmt.Sprintf("`json:\"%s,omitempty\"`", p.e.h.SpecialCamel(end.Name)),
			}
		}
	}
}
