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

export default class NewTaskModal extends React.Component {
  constructor(props) {
    super(props)
  }

  render() {
    return (
      <Modal
          isOpen={this.props.isTaskModalOpen}
          onRequestClose={this.props.closeNewTaskModal}
          style={customStyles}
      >
        <div className="task-form">
          <form className="pure-form pure-form-stacked">
            <fieldset>
              <legend>Create Task</legend>
              <label htmlFor="title">Title</label>
              <input id="title" name="title" type="text" value={this.props.newTask.title} onChange={this.props.updateNewTaskTitle} placeholder="Task Name" className="form-control" />
              <label htmlFor="description">Description</label>
              <textarea id="description" name="description" value={this.props.newTask.description} onChange={this.props.updateNewTaskDescription} placeholder="Task description" className="form-control" />
              <div className="form-action">
                <button onClick={e => this.props.fetchCreateTask(this.props.projectID, this.props.selectedList.ID,  this.props.newTask.title, this.props.newTask.description)} className="pure-button pure-button-primary" type="button">Create Task</button>
              </div>
            </fieldset>
          </form>
        </div>
      </Modal>
    )
  }
}
