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
    const { isListModalOpen, newList, lists, project, isTaskModalOpen, newTask, selectedList, isListEditModalOpen, taskDraggingFrom, taskDraggingTo, error } = this.props.ListReducer
    return (
      <div id="lists">
        <div className="flash flash-error">{error}</div>
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
                <label htmlFor="color">Color</label>
                <input id="color" name="color" type="text" value={newList.color} onChange={this.props.updateNewListColor} className="form-control" />
                <div className="form-action">
                  <button onClick={e => this.props.fetchCreateList(this.props.params.projectId, newList.title, newList.color)} className="pure-button pure-button-primary" type="button">Create List</button>
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
                  <button onClick={e => this.props.fetchCreateTask(this.props.params.projectId, selectedList.Id,  newTask.title)} className="pure-button pure-button-primary" type="button">Create Task</button>
                </div>
              </fieldset>
            </form>
          </div>
        </Modal>
        <Modal
          isOpen={isListEditModalOpen}
          onRequestClose={this.props.closeEditListModal}
          style={customStyles}
        >
          <div className="list-form">
            <form className="pure-form pure-form-stacked">
              <fieldset>
                <legend>Edit List</legend>
                <label htmlFor="title">Title</label>
                <input id="title" name="title" type="text" value={selectedList !=null ? selectedList.Title.String : ''} onChange={this.props.updateSelectedListTitle} className="form-control" />
                <label htmlFor="color">Color</label>
                <input id="color" name="color" type="text" value={selectedList !=null ? selectedList.Color.String : ''} onChange={this.props.updateSelectedListColor} className="form-control" />
                <div className="form-action">
                  <button onClick={e => this.props.fetchUpdateList(this.props.params.projectId, selectedList)} className="pure-button pure-button-primary" type="button">Update List</button>
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
              <div className="fascia-list" data-dropped-depth="0" data-id={list.Id} onDragOver={this.props.taskDragOver} onDrop={e=> this.props.taskDrop(project.Id, taskDraggingFrom, taskDraggingTo)} onDragLeave={this.props.taskDragLeave}>
                <div className="fascia-list-menu" data-dropped-depth="1"><i className="fa fa-pencil" onClick={e => this.props.openEditListModal(list)} data-dropped-depth="2"></i></div>
                <span className="list-title" data-dropped-depth="1">{list.Title}</span>
                <ul className="fascia-task" data-dropped-depth="1">
                  {list.ListTasks.map(function(task, index) {
                    if (task.draggedOn) {
                      return <li className="arrow"></li>
                    } else {
                      return <li style={{"border-left": `solid 6px #${list.Color.String}`}} className="task" draggable="true" data-dropped-depth="2" data-id={task.Id} onDragStart={this.props.taskDragStart}>{task.Title.String}</li>
                    }
                  }, this)}
                  <li className="new-task" data-dropped-depth="2" style={{"border-left": `solid 6px #${list.Color.String}`}} onClick={e => this.props.openNewTaskModal(list)}>
                    <i className="fa fa-plus" data-dropped-depth="3"></i>
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
