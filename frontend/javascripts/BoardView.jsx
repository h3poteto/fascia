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

var BoardView = function(component, text, projects, modal) {
  return (
    <div id="projects">
      <Modal
        isOpen={modal}
        onRequestClose={component.closeModal}
        style={customStyles}
      >
        <div className="project-form">
          <form className="pure-form pure-form-stacked">
            <fieldset>
              <legend>Create Project</legend>
              <label htmlfor="title">Title</label>
              <input id="title" name="title" type="text" value={text} onChange={component.updateNewText} placeholder="Project Name" className="form-control" />
              <div className="form-action">
                <button onClick={component.createProject} className="pure-button pure-button-primary" type="button">CreateProject</button>
              </div>
            </fieldset>
          </form>
        </div>
      </Modal>
      <div className="items">
        {projects.map(function(item, index) {
          return <div className="fascia-card pure-button button-secondary" data-id={item.Id}>
          <span className="card-title">{item.Title}</span>
          description
          </div>;
         }, component)}
          <button onClick={component.newProject} className="pure-button button-large fascia-new-project" type="button">New</button>
      </div>
    </div>
  );

 }

export default BoardView;
