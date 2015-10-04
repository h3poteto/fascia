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

var BoardView = function(component, text, projects, repositories, modal, selectRepository) {
  return (
    <div id="projects">
      <Modal
        isOpen={modal}
        onRequestClose={component.closeModal.bind(component)}
        style={customStyles}
      >
        <div className="project-form">
          <form className="pure-form pure-form-stacked">
            <fieldset>
              <legend>Create Project</legend>
              <label htmlFor="title">Title</label>
              <input id="title" name="title" type="text" value={text} onChange={component.updateNewText.bind(component)} placeholder="Project Name" className="form-control" />
              <label htmlFor="repositories">GitHub</label>
              <select id="repositories" name="repositories" onChange={component.changeSelectRepository.bind(component)} className="form-control">
                <option value="0">--</option>
                {repositories.map(function(repo, index) {
                  if (repo.id == selectRepository) {
                    return <option value={repo.id} selected>{repo.full_name}</option>;
                  } else {
                    return<option value={repo.id}>{repo.full_name}</option>;
                  }
                }, component)}
              </select>
              <div className="form-action">
                <button onClick={component.createProject.bind(component)} className="pure-button pure-button-primary" type="button">CreateProject</button>
              </div>
            </fieldset>
          </form>
        </div>
      </Modal>
      <div className="items">
        {projects.map(function(item, index) {
          return (
            <div className="fascia-card pure-button button-secondary" data-id={item.Id}>
            <span className="card-title">{item.Title}</span>
            description
            </div>
          );
         }, component)}
            <button onClick={component.newProject.bind(component)} className="pure-button button-large fascia-new-project" type="button">New</button>
      </div>
    </div>
  );

 }

export default BoardView;
