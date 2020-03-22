import React from 'react'
import { ThunkDispatch } from 'redux-thunk'
import { RouteComponentProps } from 'react-router-dom'

import Actions, { getLists } from '@/actions/projects/lists'
import { RootStore } from '@/reducers/index'
import List from './lists/list.tsx'
import styles from './lists.scss'

type Props = {
  dispatch: ThunkDispatch<any, any, Actions>
} & RootStore &
  RouteComponentProps<{ project_id: string }>

class ListsComponent extends React.Component<Props> {
  componentDidMount() {
    const id = this.props.match.params.project_id
    this.props.dispatch(getLists(parseInt(id)))
  }

  render() {
    const lists = this.props.lists.lists
    return (
      <div className={styles.lists}>
        {lists.map(l => (
          <List key={l.id} list={l} />
        ))}
        <div>{this.props.children}</div>
      </div>
    )
  }
}

export default ListsComponent
