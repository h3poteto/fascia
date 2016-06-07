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

export default class ListView extends React.Component {
  constructor(props) {
    super(props)
  }

  componentWillMount() {
    this.props.fetchLists(this.props.params.projectID)
    this.props.fetchProject(this.props.params.projectID)
    this.props.fetchListOptions()
  }

  componentWillReceiveProps(nextProps) {
    if (this.props.ListReducer.error != null || nextProps.ListReducer.error != null) {
      setTimeout(() => {
        this.props.closeFlash()
      }, 3000)
    }
  }

  componentDidMount() {
    let maxHeight = window.innerHeight * 0.7;
    let stylesheet = document.styleSheets.item(4);
    var idx = document.styleSheets[4].cssRules.length;
    stylesheet.insertRule("#lists .fascia-list { max-height: " + maxHeight + "px; }", idx);

  }

  wholeLoading(isLoading) {
    if (isLoading) {
      return (
        <div className="whole-loading">
          <div className="whole-circle-wrapper">
            <div className="whole-circle-body">
              <div className="whole-spinner"></div>
            </div>
          </div>
        </div>
      )
    }
  }

  listLoading() {
    return (
      <div className="list-loading">
        <div className="list-circle-wrapper">
          <div className="list-circle-body">
            <div className="list-spinner"></div>
          </div>
        </div>
      </div>
    )
  }

  listAction(project, listOptions, selectedList, selectedListOption) {
    if (project == null || project.RepositoryID == undefined || project.RepositoryID == null || project.RepositoryID == 0) {
      return null
    } else {
      return (
        <div>
          <label htmlFor="action">action</label>
          <select id="action" name="action" type="text" onChange={this.props.changeSelectedListOption} className="form-control" value={selectedListOption ? selectedListOption.ID : (selectedList ? selectedList.ListOptionID : 0)}>
            <option value="0">nothing</option>
            {listOptions.map(function(option, index) {
               return <option key={index} value={option.ID}>{option.Action}</option>
             }, this)}
          </select>
        </div>
      )
    }
  }

  projectOperations(project, selectedProject) {
    if (project == null || project.RepositoryID == undefined || project.RepositoryID == null || project.RepositoryID == 0) {
      return <span></span>
    } else {
      return (
        <span>
          <span className={project.ShowPullRequests ? "pull-request-select select" : "pull-request-select"} onClick={e => this.props.showPullRequests(this.props.params.projectID, selectedProject.ShowIssues, selectedProject.ShowPullRequests)}><i className="octicon octicon-git-pull-request"></i></span>
          <span className={project.ShowIssues ? "pull-request-select select" : "pull-request-select"} onClick={e => this.props.showIssues(this.props.params.projectID, selectedProject.ShowIssues, selectedProject.ShowPullRequests)}><i className="octicon octicon-issue-opened"></i></span>
          <i className="fa fa-repeat" onClick={e => this.props.fetchProjectGithub(this.props.params.projectID)}></i>
        </span>
      )
    }
  }

  webhookButton(project) {
    if (project == null || project.RepositoryID == undefined || project.RepositoryID == null || project.RepositoryID == 0) {
      return null
    } else {
      return (
        <button onClick={e => this.props.createWebhook(this.props.params.projectID)} className="pure-button button-secondary" type="button">Update Webhook</button>
      )
    }
  }

  listEditButton(list) {
    if (list.IsInitList) {
      return null
    } else {
      return <i className="fa fa-pencil" onClick={e => this.props.openEditListModal(list)} data-dropped-depth="2"></i>
    }
  }

  listItem(index, list, project, taskDraggingFrom, taskDraggingTo) {
    if (list.IsHidden) {
      return (
        <div key={index} className="fascia-list fascia-list-hidden" data-dropped-depth="0" data-id={list.ID} onDragOver={this.props.taskDragOver} onDrop={e=> this.props.taskDrop(project.ID, taskDraggingFrom, taskDraggingTo)} onDragLeave={this.props.taskDragLeave}>
          <div className="fascia-list-menu" data-dropped-depth="1">
            <i className="fa fa-eye" onClick={e => this.props.displayList(project.ID, list.ID)} data-dropped-depth="2"></i>
            {this.listEditButton(list)}
          </div>
          <span className="list-title" data-dropped-depth="1">{list.Title}</span>
          {list.isLoading != undefined && list.isLoading ? this.listLoading() : ''}
        </div>
      )
    } else {
      return (
        <div key={index} className="fascia-list" data-dropped-depth="0" data-id={list.ID} onDragOver={this.props.taskDragOver} onDrop={e=> this.props.taskDrop(project.ID, taskDraggingFrom, taskDraggingTo)} onDragLeave={this.props.taskDragLeave}>
          <div className="fascia-list-menu" data-dropped-depth="1">
            <i className="fa fa-eye-slash" onClick={e => this.props.hideList(project.ID, list.ID)} data-dropped-depth="2"></i>
            {this.listEditButton(list)}
          </div>
          <span className="list-title" data-dropped-depth="1">{list.Title}</span>
          <ul className="fascia-task" data-dropped-depth="1">
            {list.ListTasks.map(function(task, index) {
               if (task.draggedOn) {
                 return <li key={index} className="arrow"></li>
               } else if(project != null && project.ShowIssues && !task.PullRequest || project != null && project.ShowPullRequests && task.PullRequest) {
                 return <li key={index} style={{"borderLeft": `solid 6px #${list.Color}`}} className="task" draggable="true" data-dropped-depth="2" data-id={task.ID} onDragStart={this.props.taskDragStart}>{task.Title}</li>
               }
             }, this)}
            <li className="new-task" data-dropped-depth="2" style={{"borderLeft": `solid 6px #${list.Color}`}} onClick={e => this.props.openNewTaskModal(list)}>
              <i className="fa fa-plus" data-dropped-depth="3"></i>
            </li>
          </ul>
          {list.isLoading != undefined && list.isLoading ? this.listLoading() : ''}
        </div>
      )
    }
  }

  render() {
    const { isLoading, isListModalOpen, newList, lists, listOptions, noneList, project, isTaskModalOpen, newTask, selectedList, selectedListOption, isListEditModalOpen, isProjectEditModalOpen, taskDraggingFrom, taskDraggingTo, selectedProject, error } = this.props.ListReducer

    var flash;
    if (error != null) {
      flash = <div className="flash flash-error">{error}</div>;
    }

    return (
      <div id="lists">
        {this.wholeLoading(isLoading)}
        {flash}
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
                  <button onClick={e => this.props.fetchCreateList(this.props.params.projectID, newList.title, newList.color)} className="pure-button pure-button-primary" type="button">Create List</button>
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
                <label htmlFor="description">Description</label>
                <textarea id="description" name="description" value={newTask.description} onChange={this.props.updateNewTaskDescription} placeholder="Task description" className="form-control" />
                <div className="form-action">
                  <button onClick={e => this.props.fetchCreateTask(this.props.params.projectID, selectedList.ID,  newTask.title, newTask.description)} className="pure-button pure-button-primary" type="button">Create Task</button>
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
                <input id="title" name="title" type="text" value={selectedList !=null ? selectedList.Title : ''} onChange={this.props.updateSelectedListTitle} className="form-control" />
                <label htmlFor="color">Color</label>
                <input id="color" name="color" type="text" value={selectedList !=null ? selectedList.Color : ''} onChange={this.props.updateSelectedListColor} className="form-control" />
                {this.listAction(project, listOptions, selectedList, selectedListOption)}
                <div className="form-action">
                  <button onClick={e => this.props.fetchUpdateList(this.props.params.projectID, selectedList, selectedListOption)} className="pure-button pure-button-primary" type="button">Update List</button>
                </div>
              </fieldset>
            </form>
          </div>
        </Modal>
        <Modal
          isOpen={isProjectEditModalOpen}
          onRequestClose={this.props.closeEditProjectModal}
          style={customStyles}
        >
          <div className="project-form">
            <form className="pure-form pure-form-stacked">
              <fieldset>
                <legend>Edit Project</legend>
                <label htmlFor="title">Title</label>
                <input id="title" name="title" type="text" value={selectedProject.Title} onChange={this.props.updateEditProjectTitle} className="form-control" />
                <label htmlFor="description">Description</label>
                <textarea id="description" name="description" value={selectedProject.Description} onChange={this.props.updateEditProjectDescription} className="form-control" />
                <div className="form-action">
                  {this.webhookButton(project)}&nbsp;
                  <button onClick={e => this.props.fetchUpdateProject(this.props.params.projectID, selectedProject)} className="pure-button pure-button-primary" type="button">Update Project</button>
                </div>
              </fieldset>
            </form>
          </div>
        </Modal>
        <div className="title-wrapper">
          <div className="project-operation">
            {this.projectOperations(project, selectedProject)}
          </div>
          <h3 className="project-title">{project != null ? project.Title : ''}<span className="fascia-project-menu" onClick={e => this.props.openEditProjectModal(project)}><i className="fa fa-pencil"></i></span></h3>
        </div>
        <div className="items">
          {lists.map(function(list, index) {
            return this.listItem(index, list, project, taskDraggingFrom, taskDraggingTo)
           }, this)}
           <button onClick={this.props.openNewListModal} className="pure-button button-large fascia-new-list pure-button-primary" type="button">New</button>
           <div className="clearfix"></div>
        </div>
        <div className="none-list-tasks" data-dropped-depth="0" data-id={noneList.ID} onDragOver={this.props.taskDragOver} onDrop={e => this.props.taskDrop(project.ID, taskDraggingFrom, taskDraggingTo)} onDragLeave={this.props.taskDragLeave}>
          <ul className="fascia-none-list-tasks" data-dropped-depth="1">
            {noneList.ListTasks.map(function(task, index) {
               if (task.draggedOn) {
                 return <li key={index} className="arrow"></li>
               } else if( project != null && project.ShowIssues && !task.PullRequest || project != null && project.ShowPullRequests && task.PullRequest) {
                 return <li key={index} className="button-green task" draggable="true" data-dropped-depth="2" data-id={task.ID} onDragStart={this.props.taskDragStart}>{task.Title}</li>
               }
             }, this)}
            <li onClick={e => this.props.openNewTaskModal(noneList)} className="new-task pure-button button-blue" data-dropped-depth="2">
              <i className="fa fa-plus" data-dropped-depth="3"></i>
            </li>
          </ul>
        </div>
      </div>
    );
  }
}
