import React from 'react'

import { Task } from '@/entities/task'
import styles from './task.scss'

type Props = {
  task: Task
  color: string
}

const task: React.FC<Props> = props => {
  const border = {
    'borderLeft': `6px solid #${props.color}`
  } as React.CSSProperties

  return (
    <div className={styles.task} style={border}>
      {props.task.title}
    </div>
  )
}

export default task
