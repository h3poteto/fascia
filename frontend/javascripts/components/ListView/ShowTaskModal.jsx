import React from 'react'
import Modal from 'react-modal'

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

  render() {
    return (
      <Modal
          isOpen={this.props.isShowTaskModalOpen}
          onRequestClose={this.props.closeShowTaskModal}
          style={customStyles}
      >
        <div className="task-detail">
          <div className="task-title">
            <span className="octicon octicon-mark-github task-icon"></span>
            {this.props.task.Title}
            {this.issueNumber(this.props.task)}
          </div>
          <div className="task-description">
            {this.props.task.Description}
          </div>
        </div>
      </Modal>
    )
  }
}
