import React from 'react'
import { Button } from 'react-bootstrap'
import { ThunkDispatch } from 'redux-thunk'
import { Link, RouteComponentProps } from 'react-router-dom'

import Actions, { getLists, getProject, openDelete, closeDelete, openNewList, closeNewList, openEditProject, closeEditProject } from '@/actions/projects/lists'
import { RootStore } from '@/reducers/index'
import List from './lists/list.tsx'
import styles from './lists.scss'
import Delete from './lists/delete.tsx'
import New from './lists/new.tsx'
import Task from './lists/list/task.tsx'
import EditProject from './lists/editProject.tsx'

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

    const id = parseInt(this.props.match.params.project_id)

    const lists = this.props.lists.lists
    const project = this.props.lists.project
    const noneList = this.props.lists.noneList
    return (
      <div>
        <div className={styles.title}>
          <h3>{ project ? project.title : '' }</h3>
          <span className="mr-2" onClick={openEditProjectModal}><i title="Edit project" className="fa fa-pencil"></i></span>
          <span onClick={openDeleteModal}><i title="Delete project" className="fa fa-trash"></i></span>
        </div>
        <div className={styles.backboard}>
          <div className={styles.lists}>
            {lists.map(l => (
              <List key={l.id} list={l} />
            ))}
            <Button className={styles.newButton} onClick={openNewListModal}>New</Button>
          </div>
          <div className={styles.noneList}>
            {noneList && noneList.tasks.map(t => (
              <Link key={t.id} to={`/projects/${noneList.project_id}/lists/${noneList.id}/tasks/${t.id}`}>
                <Task key={t.id} task={t} color="218838" />
              </Link>
            ))}
            {noneList && (
              <Link to={`/projects/${noneList.project_id}/lists/${noneList.id}/tasks/new`}>
                <Button style={{ width: '100%' }} variant="outline-info"><i className="fa fa-plus"></i></Button>
              </Link>
            )}
          </div>
        </div>
        <New open={this.props.lists.newListModal} close={closeNewListModal} color={this.props.lists.defaultColor} projectID={id} dispatch={this.props.dispatch} />
        <Delete open={this.props.lists.deleteModal} project={this.props.lists.project} close={closeDeleteModal} dispatch={this.props.dispatch} />
        <EditProject open={this.props.lists.editProjectModal} close={closeEditProjectModal} project={project} dispatch={this.props.dispatch} />
      </div>
    )
  }
}

export default ListsComponent
