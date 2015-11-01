import React from 'react';
import Modal from 'react-modal';

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
};

export default class ListView extends React.Component {
  constructor(props) {
    super(props);
  }

  componentWillMount() {
    this.props.fetchLists(this.props.params.projectId);
    this.props.fetchProject(this.props.params.projectId);
  }

  render() {
    const { isListModalOpen, newList, lists, project, isTaskModalOpen, newTask, selectedListId } = this.props.ListReducer
    return (
      <div id="lists">
        <Modal
          isOpen={isListModalOpen}
          onRequestClose={this.props.closeNewListModal}
          style={customStyles}
        >
          <div className="list-form">
            <form className="pure-form pure-form-stacked">
              <fieldset>
                <legend>Create List</legend>
                <label htmlFor="title">Title</label>
                <input id="title" name="title" type="text" value={newList.title} onChange={this.props.updateNewListTitle} placeholder="List Name" className="form-control" />
                <div className="form-action">
                  <button onClick={e => this.props.fetchCreateList(this.props.params.projectId, newList.title)} className="pure-button pure-button-primary" type="button">Create List</button>
                </div>
              </fieldset>
            </form>
          </div>
        </Modal>
        <Modal
          isOpen={isTaskModalOpen}
          onRequestClose={this.props.closeNewTaskModal}
          style={customStyles}
        >
          <div className="task-form">
            <form className="pure-form pure-form-stacked">
              <fieldset>
                <legend>Create Task</legend>
                <label htmlFor="title">Title</label>
                <input id="title" name="title" type="text" value={newTask.title} onChange={this.props.updateNewTaskTitle} placeholder="Task Name" className="form-control" />
                <div className="form-action">
                  <button onClick={e => this.props.fetchCreateTask(this.props.params.projectId, selectedListId,  newTask.title)} className="pure-button pure-button-primary" type="button">Create Task</button>
                </div>
              </fieldset>
            </form>
          </div>
        </Modal>
        <div className="title-wrapper">
          <h3 className="project-title">{project != null ? project.Title : ''}</h3>
        </div>
        <div className="items">
          {lists.map(function(list, index) {
            return (
              <div className="fascia-list" data-id={list.Id}>
                <span className="list-title">{list.Title}</span>
                <ul className="fascia-task">
                  {list.ListTasks.map(function(task, index) {
                    return <li className="task">{task.Title.String}</li>
                  }, this)}
                  <li className="new-task" onClick={e => this.props.openNewTaskModal(list.Id)}>
                    <i className="fa fa-plus"></i>
                  </li>
                </ul>
              </div>
            );
           }, this)}
              <button onClick={this.props.openNewListModal} className="pure-button button-large fascia-new-list pure-button-primary" type="button">New</button>
        </div>
      </div>
    );
  }
}
