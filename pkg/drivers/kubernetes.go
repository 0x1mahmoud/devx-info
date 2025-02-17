package drivers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/encoding/yaml"
	log "github.com/sirupsen/logrus"
	"github.com/stakpak/devx/pkg/stack"
	"github.com/stakpak/devx/pkg/stackbuilder"
	"github.com/stakpak/devx/pkg/utils"
)

type KubernetesDriver struct {
	Config stackbuilder.DriverConfig
}

func (d *KubernetesDriver) match(resource cue.Value) bool {
	driverName, _ := resource.LookupPath(cue.ParsePath("$metadata.labels.driver")).String()
	return driverName == "kubernetes"
}

type Manifest struct {
	Name string
	Data []byte
}

func (d *KubernetesDriver) ApplyAll(stack *stack.Stack, stdout bool) error {
	manifests := []Manifest{}

	for _, componentId := range stack.GetTasks() {
		component, _ := stack.GetComponent(componentId)

		resourceIter, _ := component.LookupPath(cue.ParsePath("$resources")).Fields()
		for resourceIter.Next() {
			if d.match(resourceIter.Value()) {
				resource, err := utils.RemoveMeta(resourceIter.Value())
				if err != nil {
					return err
				}

				kind := resource.LookupPath(cue.ParsePath("kind"))
				if kind.Err() != nil {
					return kind.Err()
				}

				kindString, err := kind.String()
				if err != nil {
					return err
				}

				name := resource.LookupPath(cue.ParsePath("metadata.name"))
				if kind.Err() != nil {
					return kind.Err()
				}

				nameString, err := name.String()
				if err != nil {
					return err
				}

				data, err := yaml.Encode(resource)
				if err != nil {
					return err
				}

				manifests = append(manifests, Manifest{
					Name: fmt.Sprintf("%s-%s.yml", nameString, strings.ToLower(kindString)),
					Data: data,
				})
			}
		}
	}

	if len(manifests) == 0 {
		return nil
	}

	if stdout {
		for _, m := range manifests {
			if _, err := os.Stdout.Write([]byte("---\n")); err != nil {
				return err
			}
			if _, err := os.Stdout.Write(m.Data); err != nil {
				return err
			}
		}
		_, err := os.Stdout.Write([]byte("\n"))
		return err
	}

	if _, err := os.Stat(d.Config.Output.Dir); os.IsNotExist(err) {
		os.MkdirAll(d.Config.Output.Dir, 0700)
	}

	if d.Config.Output.File == "" {
		for _, m := range manifests {
			filePath := filepath.Join(d.Config.Output.Dir, m.Name)
			os.WriteFile(filePath, m.Data, 0700)
		}
		log.Infof("[kubernetes] applied resources to \"%s/*.yml\"", d.Config.Output.Dir)
		return nil
	}

	filePath := filepath.Join(d.Config.Output.Dir, d.Config.Output.File)
	_, err := os.Stat(filePath)
	if !os.IsNotExist(err) {
		os.Remove(filePath)
	}
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0700)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for _, m := range manifests {
		file.Write([]byte("---\n"))
		file.Write(m.Data)
	}
	log.Infof("[kubernetes] applied resources to \"%s\"", filePath)

	return nil
}
