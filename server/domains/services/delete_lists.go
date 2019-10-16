package services

import (
	"github.com/h3poteto/fascia/server/domains/list"
	"github.com/h3poteto/fascia/server/domains/project"
)

// DeleteLists delete all lists related a project
func DeleteLists(p *project.Project, listInfra list.Repository) error {
	lists, err := listInfra.Lists(p.ID)
	if err != nil {
		return err
	}
	for _, l := range lists {
		err := listInfra.DeleteTasks(l.ID)
		if err != nil {
			return err
		}
		err = listInfra.Delete(l.ID)
		if err != nil {
			return err
		}
	}
	noneList, err := listInfra.NoneList(p.ID)
	err = listInfra.DeleteTasks(noneList.ID)
	if err != nil {
		return err
	}
	return listInfra.Delete(noneList.ID)
}
