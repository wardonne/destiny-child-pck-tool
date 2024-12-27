package pcktool

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"pcktool/object"
)

func GenerateLive2D(pck *object.Package) (*object.Live2DModel, error) {
	var model = new(object.Live2DModel)
	// find model.json file
	for _, entry := range pck.Entries {
		if entry.Ext != "json" {
			continue
		}
		var testModel = new(object.Live2DModel)
		if err := json.Unmarshal(entry.Content, testModel); err != nil {
			panic(err)
		}
		if testModel.Version == nil || testModel.Model == nil || testModel.Textures == nil {
			continue
		}
		model = testModel
		entry.Filename = "model.json"
		entry.Ext = ""
		break
	}
	if model == nil {
		return nil, fmt.Errorf("Live2D model not found")
	}
	// find .dat file
	for _, entry := range pck.Entries {
		if entry.Ext != "dat" {
			continue
		}
		entry.Filename = *model.Model
		entry.Ext = ""
		break
	}
	// find .png file
	var textures []*object.PackageEntry
	for _, entry := range pck.Entries {
		if entry.Ext != "png" {
			continue
		}
		textures = append(textures, entry)
	}
	for index, texture := range *model.Textures {
		textures[index].Filename = filepath.ToSlash(filepath.Join("textures", texture))
		textures[index].Ext = ""
		(*model.Textures)[index] = textures[index].Filename
	}
	// find .mtn files
	var motions []*object.PackageEntry
	for _, entry := range pck.Entries {
		if entry.Ext != "mtn" {
			continue
		}
		motions = append(motions, entry)
	}
	index := 0
	for motionName, motionArray := range *model.Motions {
		if len(motionArray) == 0 {
			continue
		}
		for i, motion := range motionArray {
			motions[index].Filename = filepath.ToSlash(filepath.Join("motions", motion.File))
			motions[index].Ext = ""
			(*model.Motions)[motionName][i].File = motions[index].Filename
		}
		index++
	}
	// find .exp.json files
	if model.Expressions != nil {
		var expressions []*object.PackageEntry
		for _, entry := range pck.Entries {
			if entry.Ext != "json" {
				continue
			}
			expressions = append(expressions, entry)
		}
		for index, expression := range *model.Expressions {
			expressions[index].Filename = filepath.ToSlash(filepath.Join("expressions", expression.File))
			expressions[index].Ext = ""
			(*model.Expressions)[index].File = expressions[index].Filename
		}
	}
	// find model.json and replace content
	for _, entry := range pck.Entries {
		if entry.Filename != "model.json" {
			continue
		}
		var err error
		entry.Content, err = json.MarshalIndent(model, "", "  ")
		if err != nil {
			return nil, err
		}
	}
	return model, nil
}
