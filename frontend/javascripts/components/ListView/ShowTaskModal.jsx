import React from 'react'
import PropTypes from 'prop-types'
import Modal from 'react-modal'
import MarkdownIt from 'markdown-it'
import MarkdownItCheckbox from 'markdown-it-checkbox'
import { Field, reduxForm } from 'redux-form'

import { RenderField, validate } from './taskForm'

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
    if (!nextProps.dirty || !nextProps.isShowTaskModalOpen) {
      this.handleInitialize(nextProps)
    }
  }

  handleInitialize(props) {
    const initData = {
      'title': props.task.Title,
      'description': props.task.Description,
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
      return 'Description'
    }
  }

  taskIcon(project) {
    if (project.RepositoryID != undefined && project.RepositoryID != null && project.RepositoryID != 0) {
      return <span className="octicon octicon-mark-github task-icon"></span>
    } else {
      return <img className="fascia-icon task-icon" src="/images/fascia-icon.png" />
    }
  }

  taskForm(project, task, isEditTaskModalVisible, handleSubmit, action, pristine, submitting, reset, changeEditMode) {
    if (isEditTaskModalVisible) {
      return (
        <div className="task-body task-form">
          <form className="pure-form pure-form-stacked" onSubmit={handleSubmit(action)}>
            <fieldset>
              <legend>Edit Task</legend>
              <label htmlFor="title">Title</label>
              <Field component={RenderField} name="title" id="title" type="text" />
              <label htmlFor="description">Description</label>
              <Field name="description" id="description" component="textarea" className="form-control" />
              <div className="form-action">
                <button type="reset" className="pure-button pure-button-default" onClick={() => changeEditMode(task)}>Cancel</button>
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
            {this.taskIcon(project)}
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

  deleteTask(task, fetchDeleteTask) {
    if (task.IssueNumber === 0) {
      return <i title="Delete task" className="fa fa-trash" onClick={fetchDeleteTask}></i>
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
      project,
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
          contentLabel="ShowTaskModal"
      >
        <div className="task-detail">
          <div className="task-controll">
            {this.deleteTask(task, fetchDeleteTask)}
            <i title="Edit task" className="fa fa-pencil" onClick={() => changeEditMode(this.props.task)}></i>
          </div>
          {this.taskForm(project, task, isEditTaskModalVisible, handleSubmit, action, pristine, submitting, reset, changeEditMode)}
        </div>
      </Modal>
    )
  }
}

ShowTaskModal.propTypes = {
  initialize: PropTypes.func.isRequired,
  handleSubmit: PropTypes.func.isRequired,
  pristine: PropTypes.bool,
  reset: PropTypes.func.isRequired,
  submitting: PropTypes.bool.isRequired,
  onRequestClose: PropTypes.func.isRequired,
  action: PropTypes.func.isRequired,
  project: PropTypes.object,
  task: PropTypes.object,
  isShowTaskModalOpen: PropTypes.bool.isRequired,
  isEditTaskModalVisible: PropTypes.bool,
  fetchDeleteTask: PropTypes.func.isRequired,
  changeEditMode: PropTypes.func.isRequired,
  dirty: PropTypes.object,
}

export default reduxForm({
  form: 'edit-task-form',
  validate,
})(ShowTaskModal)
