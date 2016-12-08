import React from 'react'
import Modal from 'react-modal'
import MarkdownIt from 'markdown-it'
import MarkdownItCheckbox from 'markdown-it-checkbox'
import { Field, reduxForm } from 'redux-form'

const customStyles = {
  overlay : {
    position          : 'fixed',
    top               : 0,
    left              : 0,
    right             : 0,
    bottom            : 0,
    backgroundColor   : 'rgba(255, 255, 255, 0.5)'
  },
  content : {
    position : 'fixed',
    top : '50%',
    left : '50%',
    right : 'auto',
    bottom : 'auto',
    overflow : 'auto',
    marginRight : '-50%',
    transform : 'translate(-50%, -50%)'
  }
}

class ShowTaskModal extends React.Component {
  componentWillReceiveProps(nextProps) {
    // modalをcloseするタイミングでは初期化しておかないと，別のtaskを選択したときに，現在の編集分が残っている可能性がある
    console.log(nextProps.dirty)
    if (!nextProps.dirty || !nextProps.isShowTaskModalOpen) {
      this.handleInitialize(nextProps)
    }
  }

  handleInitialize(props) {
    const initData = {
      "title": props.task.Title,
      "description": props.task.Description,
    }

    this.props.initialize(initData)
  }

  issueNumber(task) {
    if (task.IssueNumber == 0) {
      return <span></span>
    } else {
      return (
        <a href={task.HTMLURL} target="_blank">
          <span className="task-issue-number">
            #{task.IssueNumber}
          </span>
        </a>
      )
    }
  }

  markdownDescription(task) {
    if (task.Description.length > 0) {
      let rawMarkup = MarkdownIt({
        html: true,
        linkify: true,
        breaks: true,
        typographer: true
      })
        .use(MarkdownItCheckbox)
        .render(task.Description)
      return <span dangerouslySetInnerHTML={{__html: rawMarkup}} />
    } else {
      return "Description"
    }
  }

  taskForm(projectID, task, isEditTaskModalVisible, handleSubmit, action, pristine, submitting, reset) {
    if (isEditTaskModalVisible) {
      return (
        <div className="task-body task-form">
          <form className="pure-form pure-form-stacked" onSubmit={handleSubmit((values) => { action(projectID, task.ListID, task.ID, values) })}>
            <fieldset>
              <legend>Edit Task</legend>
              <label htmlFor="title">Title</label>
              <Field name="title" id="title" component="input" type="text" className="form-control" />
              <label htmlFor="description">Description</label>
              <Field name="description" id="description" component="textarea" className="form-control" />
              <div className="form-action">
                <button type="reset" className="pure-button pure-button-default" disabled={pristine || submitting} onClick={reset}>Reset</button>
                <button type="submit" className="pure-button pure-button-primary" disabled={pristine || submitting}>Update Task</button>
              </div>
            </fieldset>
          </form>
        </div>
      )
    } else {
      return (
        <div className="task-body">
          <div className="task-title">
            <span className="octicon octicon-mark-github task-icon"></span>
            {task.Title}
            {this.issueNumber(task)}
          </div>
          <div className="task-description">
            {this.markdownDescription(task)}
          </div>
        </div>
      )
    }
  }

  deleteTask(projectID, task, fetchDeleteTask) {
    if (task.IssueNumber === 0) {
      return <i title="Delete task" className="fa fa-trash" onClick={() => fetchDeleteTask(projectID, task.ListID, task.ID)}></i>
    } else {
      return
    }
  }

  render() {
    const {
      handleSubmit,
      pristine,
      reset,
      submitting,
      onRequestClose,
      action,
      projectID,
      task,
      fetchDeleteTask,
      isEditTaskModalVisible,
      isShowTaskModalOpen,
      changeEditMode,
    } = this.props
    return (
      <Modal
          isOpen={isShowTaskModalOpen}
          onRequestClose={onRequestClose}
          style={customStyles}
      >
        <div className="task-detail">
          <div className="task-controll">
            {this.deleteTask(projectID, task, fetchDeleteTask)}
            <i title="Edit task" className="fa fa-pencil" onClick={() => changeEditMode(this.props.task)}></i>
          </div>
          {this.taskForm(projectID, task, isEditTaskModalVisible, handleSubmit, action, pristine, submitting, reset)}
        </div>
      </Modal>
    )
  }
}

export default reduxForm({
  form: 'edit-task-form',
})(ShowTaskModal)
