import React from 'react'
import { Button } from 'react-bootstrap'
import { ThunkDispatch } from 'redux-thunk'
import { Link, RouteComponentProps } from 'react-router-dom'
import { DragDropContext, Droppable, Draggable, DropResult } from 'react-beautiful-dnd'

import Actions, {
  getLists,
  openDelete,
  closeDelete,
  openNewList,
  closeNewList,
  openEditProject,
  closeEditProject,
  syncGithub,
  moveTask
} from '@/actions/projects/lists'
import { getProject } from '@/actions/projects/show'
import { RootStore } from '@/reducers/index'
import ListComponent from './lists/list.tsx'
import styles from './lists.scss'
import Delete from './lists/delete.tsx'
import New from './lists/new.tsx'
import TaskComponent from './lists/list/task.tsx'
import EditProject from './lists/editProject.tsx'
import { List } from '@/entities/list'

type Props = {
  dispatch: ThunkDispatch<any, any, Actions>
} & RootStore &
  RouteComponentProps<{ project_id: string }>

class ListsComponent extends React.Component<Props> {
  componentDidMount() {
    const id = this.props.match.params.project_id
    this.props.dispatch(getLists(parseInt(id)))
    this.props.dispatch(getProject(parseInt(id)))
  }

  operations() {
    const syncRepository = () => {
      const id = parseInt(this.props.match.params.project_id)
      this.props.dispatch(syncGithub(id))
    }

    if (this.props.lists.project?.repositoryID) {
      return (
        <div className="float-end pe-5 pt-2">
          <span onClick={syncRepository}>
            <i title="Sync GitHub issues" className="fa fa-repeat"></i>
          </span>
        </div>
      )
    } else {
      return null
    }
  }

  render() {
    const openDeleteModal = () => {
      this.props.dispatch(openDelete())
    }

    const closeDeleteModal = () => {
      this.props.dispatch(closeDelete())
    }

    const openNewListModal = () => {
      this.props.dispatch(openNewList())
    }

    const closeNewListModal = () => {
      this.props.dispatch(closeNewList())
    }

    const openEditProjectModal = () => {
      this.props.dispatch(openEditProject())
    }

    const closeEditProjectModal = () => {
      this.props.dispatch(closeEditProject())
    }

    const onDragEnd = (result: DropResult) => {
      const { source, destination, draggableId } = result

      if (!destination) {
        return
      }

      const projectID = parseInt(this.props.match.params.project_id)
      const fromListID = parseInt(source.droppableId)
      const toListID = parseInt(destination.droppableId)
      const taskID = parseInt(draggableId)
      let list: List | null | undefined = this.props.lists.lists.find((e) => e.id === toListID)
      if (!list) {
        list = this.props.lists.noneList
      }
      let prevToTaskID: number | null = null
      if (list && list.tasks[destination.index]) {
        prevToTaskID = list.tasks[destination.index].id
      }

      this.props.dispatch(moveTask(projectID, fromListID, toListID, taskID, prevToTaskID))
    }

    const id = parseInt(this.props.match.params.project_id)

    const lists = this.props.lists.lists
    const project = this.props.lists.project
    const noneList = this.props.lists.noneList
    return (
      <div>
        <div className={styles.title}>
          <h3>{project ? project.title : ''}</h3>
          <span className="me-2" onClick={openEditProjectModal}>
            <i title="Edit project" className="fa fa-pencil"></i>
          </span>
          <span className="me-2" onClick={openDeleteModal}>
            <i title="Delete project" className="fa fa-trash"></i>
          </span>
          {this.operations()}
          <div className="clearfix"></div>
        </div>
        <div className={styles.backboard}>
          <DragDropContext onDragEnd={onDragEnd}>
            <div className={styles.lists}>
              {lists.map((l) => (
                <Droppable key={l.id} droppableId={`${l.id}`}>
                  {(provided) => (
                    <div ref={provided.innerRef}>
                      <ListComponent list={l} />
                      {provided.placeholder}
                    </div>
                  )}
                </Droppable>
              ))}
              <Button className={styles.newButton} onClick={openNewListModal}>
                New
              </Button>
            </div>
            <div className={styles.noneList}>
              <Droppable droppableId={`${noneList?.id}`}>
                {(provided) => (
                  <div ref={provided.innerRef}>
                    {noneList &&
                      noneList.tasks.map((t, index) => (
                        <Draggable key={t.id} draggableId={`${t.id}`} index={index}>
                          {(provided) => (
                            <div ref={provided.innerRef} {...provided.draggableProps} {...provided.dragHandleProps}>
                              <Link key={t.id} to={`/projects/${noneList.project_id}/lists/${noneList.id}/tasks/${t.id}`}>
                                <TaskComponent key={t.id} task={t} color="218838" />
                              </Link>
                            </div>
                          )}
                        </Draggable>
                      ))}
                    {provided.placeholder}
                  </div>
                )}
              </Droppable>
              {noneList && (
                <Link to={`/projects/${noneList.project_id}/lists/${noneList.id}/tasks/new`}>
                  <Button style={{ width: '100%' }} variant="outline-info">
                    <i className="fa fa-plus"></i>
                  </Button>
                </Link>
              )}
            </div>
          </DragDropContext>
        </div>
        <New
          open={this.props.lists.newListModal}
          close={closeNewListModal}
          color={this.props.lists.defaultColor}
          projectID={id}
          dispatch={this.props.dispatch}
        />
        <Delete
          open={this.props.lists.deleteModal}
          project={this.props.lists.project}
          close={closeDeleteModal}
          dispatch={this.props.dispatch}
        />
        <EditProject
          open={this.props.lists.editProjectModal}
          close={closeEditProjectModal}
          project={project}
          dispatch={this.props.dispatch}
        />
      </div>
    )
  }
}

export default ListsComponent
