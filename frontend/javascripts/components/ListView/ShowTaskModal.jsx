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

  render() {
    return (
      <Modal
          isOpen={this.props.isShowTaskModalOpen}
          onRequestClose={this.props.closeShowTaskModal}
          style={customStyles}
      >
        <div className="task-detail">
          <div className="task-title">
            {this.props.task.Title}<span className="task-issue-number">#{this.props.task.IssueNumber}</span>
          </div>
          <div className="task-description">
            {this.props.task.Description}
          </div>
        </div>
      </Modal>
    )
  }
}
