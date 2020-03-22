import React from 'react'
import { Card } from 'react-bootstrap'

import styles from './list.scss'
import { List } from '@/actions/projects/lists'
import Task from './list/task.tsx'

type Props = {
  list: List
}

const list: React.FC<Props> = props => (
  <div className={styles.list}>
    <Card bg="light" style={{ width: '18rem' }}>
      <Card.Header>{props.list.title}</Card.Header>
      <Card.Body className={styles.tasks}>
        {props.list.tasks.map(t => (
          <Task key={t.id} task={t} color={props.list.color} />
        ))}
      </Card.Body>
    </Card>
  </div>
)

export default list
