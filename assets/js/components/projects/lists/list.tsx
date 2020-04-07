import React from 'react'
import { Card, Button } from 'react-bootstrap'
import { Link } from 'react-router-dom'

import styles from './list.scss'
import Task from './list/task.tsx'
import { List } from '@/entities/list'

type Props = {
  list: List,
}

const list: React.FC<Props> = props => (
  <div className={styles.list}>
    <Card bg="light" style={{ width: '18rem' }}>
      <Card.Header>{props.list.title}</Card.Header>
      <Card.Body className={styles.tasks}>
        {props.list.tasks.map(t => (
          <div key={t.id}>
            <Link to={`/projects/${props.list.project_id}/lists/${props.list.id}/tasks/${t.id}`}>
              <Task key={t.id} task={t} color={props.list.color} />
            </Link>
          </div>
        ))}
        <Link to={`/projects/${props.list.project_id}/lists/${props.list.id}/tasks/new`}>
          <Button style={{ width: '100%' }} variant="outline-info"><i className="fa fa-plus"></i></Button>
        </Link>
      </Card.Body>
    </Card>
  </div>
)

export default list
