import React from 'react'
import { ThunkDispatch } from 'redux-thunk'
import { RouteComponentProps } from 'react-router-dom'

import Actions, { getLists, getProject, openDelete, closeDelete } from '@/actions/projects/lists'
import { RootStore } from '@/reducers/index'
import List from './lists/list.tsx'
import styles from './lists.scss'
import Delete from './lists/delete.tsx'


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

    const lists = this.props.lists.lists
    const project = this.props.lists.project
    return (
      <div>
        <div className={styles.title}>
          <h3>{ project ? project.title : '' }</h3>
          <span onClick={openDeleteModal}><i title="Delete project" className="fa fa-trash"></i></span>
        </div>
        <div className={styles.lists}>
          {lists.map(l => (
            <List key={l.id} list={l} />
          ))}
        </div>
        <Delete open={this.props.lists.deleteModal} project={this.props.lists.project} close={closeDeleteModal} dispatch={this.props.dispatch} />
      </div>
    )
  }
}

export default ListsComponent
