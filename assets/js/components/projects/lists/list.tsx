import React from 'react'
import { Card, Button } from 'react-bootstrap'
import { Link } from 'react-router-dom'
import { Draggable } from 'react-beautiful-dnd'
import { ThunkDispatch } from 'redux-thunk'

import Actions, { hideList, displayList } from '@/actions/projects/lists'
import styles from './list.scss'
import Task from './list/task.tsx'
import { List } from '@/entities/list'

type Props = {
  list: List
  project_id: number
} & {
  dispatch: ThunkDispatch<any, any, Actions>
}

class ListComponent extends React.Component<Props> {
  listEdit(list: List) {
    if (list.is_init_list) return null
    return (
      <Link to={`/projects/${list.project_id}/lists/${list.id}/edit`} className="float-end m-0">
        <i className="fa fa-pencil"></i>
      </Link>
    )
  }

  listHide(list: List) {
    const callHideList = () => {
      this.props.dispatch(hideList(this.props.project_id, list.id))
    }
    const callDisplayList = () => {
      this.props.dispatch(displayList(this.props.project_id, list.id))
    }
    if (list.is_hidden) {
      return (
        <span role="button" className="float-end ms-2" onClick={callDisplayList}>
          <i className="fa fa-chevron-down"></i>
        </span>
      )
    } else {
      return (
        <span role="button" className="float-end ms-2" onClick={callHideList}>
          <i className="fa fa-chevron-up"></i>
        </span>
      )
    }
  }

  displayData(list: List) {
    if (list.is_hidden) {
      return []
    } else {
      return list.tasks
    }
  }

  render() {
    const { list } = this.props
    return (
      <div className={styles.list}>
        <Card bg="light" style={{ width: '18rem' }}>
          <Card.Header>
            {list.title}
            {this.listHide(list)}
            {this.listEdit(list)}
            <div className="clearfix"></div>
          </Card.Header>
          <Card.Body className={styles.tasks}>
            {this.displayData(list).map((t, index) => (
              <Draggable key={t.id} draggableId={`${t.id}`} index={index}>
                {(provided) => (
                  <div ref={provided.innerRef} {...provided.draggableProps} {...provided.dragHandleProps}>
                    <Link to={`/projects/${list.project_id}/lists/${list.id}/tasks/${t.id}`} className={styles.task}>
                      <Task key={t.id} task={t} color={list.color} />
                    </Link>
                  </div>
                )}
              </Draggable>
            ))}
            <Link to={`/projects/${list.project_id}/lists/${list.id}/tasks/new`}>
              <Button style={{ width: '100%' }} variant="outline-info">
                <i className="fa fa-plus"></i>
              </Button>
            </Link>
          </Card.Body>
        </Card>
      </div>
    )
  }
}

export default ListComponent
