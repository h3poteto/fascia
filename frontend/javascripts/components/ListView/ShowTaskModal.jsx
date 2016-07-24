import React from 'react'
import Modal from 'react-modal'
import MarkdownIt from 'markdown-it'
import MarkdownItCheckbox from 'markdown-it-checkbox'

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

export default class ShowTaskModal extends React.Component {
  constructor(props) {
    super(props)
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

  taskForm(task) {
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

  render() {
    return (
      <Modal
          isOpen={this.props.isShowTaskModalOpen}
          onRequestClose={this.props.closeShowTaskModal}
          style={customStyles}
      >
        <div className="task-detail">
          <div className="task-controll">
            <i title="Edit task" className="fa fa-pencil" onClick={this.props.changeEditMode}></i>
          </div>
          {this.taskForm(this.props.task)}
        </div>
      </Modal>
    )
  }
}
